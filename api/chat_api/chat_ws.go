package chat_api

import (
	"encoding/json"
	"fmt"

	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/ctype/chat_msg"
	"github.com/LiangNing7/BlogX/models/enum/chat_msg_type"
	"github.com/LiangNing7/BlogX/models/enum/relationship_enum"
	"github.com/LiangNing7/BlogX/service/focus_service"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/LiangNing7/BlogX/utils/xss"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var UP = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var OnlineMap = map[uint]map[string]*websocket.Conn{}

type ChatRequest struct {
	RevUserID uint                  `json:"revUserID"` // 发给谁
	MsgType   chat_msg_type.MsgType `json:"msgType"`   // 1 文本 2 图片  3 md
	Msg       chat_msg.ChatMsg      `json:"msg"`       // 消息主体
}
type ChatResponse struct {
	ChatListResponse
}

func (ChatApi) ChatView(c *gin.Context) {
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil || claims == nil {
		res.FailWithMsg("请登录", c)
		return
	}

	conn, err := UP.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Errorf("ws升级失败 %s", err)
		return
	}

	userID := claims.UserID
	var user models.UserModel
	err = global.DB.Take(&user, userID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}
	addr := conn.RemoteAddr().String()
	addrMap, ok := OnlineMap[userID]
	if !ok {
		OnlineMap[userID] = map[string]*websocket.Conn{
			addr: conn,
		}
	} else {
		_, ok1 := addrMap[addr]
		if !ok1 {
			OnlineMap[userID][addr] = conn
		}
	}
	fmt.Println("进入", OnlineMap)
	for {
		// 消息类型，消息，错误
		_, p, err1 := conn.ReadMessage()
		if err1 != nil {
			// 一般是客户端断开 // websocket: close 1005 (no status)
			fmt.Println(err)
			break
		}

		var req ChatRequest
		err2 := json.Unmarshal(p, &req)
		if err2 != nil {
			res.SendConnFailWithMsg("参数错误", conn)
			continue
		}
		// 判断接收人在不在
		var revUser models.UserModel
		err = global.DB.Take(&revUser, req.RevUserID).Error
		if err != nil {
			res.SendConnFailWithMsg("接收人不存在", conn)
			continue
		}

		// 具体的消息类型做处理
		switch req.MsgType {
		case chat_msg_type.TextMsgType:
			if req.Msg.TextMsg == nil || req.Msg.TextMsg.Content == "" {
				res.SendConnFailWithMsg("文本消息内容为空", conn)
				continue
			}
		case chat_msg_type.ImageMsgType:
			if req.Msg.ImageMsg == nil || req.Msg.ImageMsg.Src == "" {
				res.SendConnFailWithMsg("图片消息内容为空", conn)
				continue
			}
		case chat_msg_type.MarkdownMsgType:
			if req.Msg.MarkdownMsg == nil || req.Msg.MarkdownMsg.Content == "" {
				res.SendConnFailWithMsg("markdown消息内容为空", conn)
				continue
			}
			// 对markdown消息做过滤
			req.Msg.MarkdownMsg.Content = xss.XSSFilter(req.Msg.MarkdownMsg.Content)
		default:
			res.SendConnFailWithMsg("不支持的消息类型", conn)
			continue
		}

		// 判断你与对方的好友关系
		// 好友就能每天聊
		// 已关注和粉丝 如果对方没有回复你，那么每天只能聊一次  对方如果没有回你，那么你只能发三条消息
		// 陌生人，如果对方开了陌生人私信，那么就能聊
		relation := focus_service.CalcUserRelationship(userID, req.RevUserID)

		fmt.Printf("用户%d %d的好友关系是：%d\n", userID, req.RevUserID, relation)
		switch relation {
		case relationship_enum.RelationStranger: // 陌生人
			var revUserMsgConf models.UserMessageConfModel
			err = global.DB.Take(&revUserMsgConf, "user_id = ?", revUser.ID).Error
			if err != nil {
				res.SendConnFailWithMsg("接收人隐私设置不存在", conn)
				continue
			}
			if !revUserMsgConf.OpenPrivateChat {
				res.SendConnFailWithMsg("对方未开始陌生人私聊", conn)
				continue
			}
		case relationship_enum.RelationFocus, relationship_enum.RelationFans: // 已关注
			// 今天对方如果没有回复你，那么你就只能发一条
			var chatList []models.ChatModel
			global.DB.Find(&chatList, "date(created_at) = date (now()) and ( (send_user_id = ? and  rev_user_id = ?) or (send_user_id = ? and  rev_user_id = ?))",
				userID, req.RevUserID, req.RevUserID, userID)

			// 我发的  对方发的
			var sendChatCount, revChatCount int
			for _, model := range chatList {
				if model.SendUserID == userID {
					sendChatCount++
				}
				if model.RevUserID == userID {
					revChatCount++
				}
			}
			fmt.Println(sendChatCount, revChatCount)
			if sendChatCount > 1 && revChatCount == 0 {
				res.SendConnFailWithMsg("对方未回复的情况下，当天只能发送一条消息", conn)
				continue
			}

		}

		// 先落库
		model := models.ChatModel{
			SendUserID: claims.UserID,
			RevUserID:  req.RevUserID,
			MsgType:    req.MsgType,
			Msg:        req.Msg,
		}
		err = global.DB.Create(&model).Error
		if err != nil {
			res.SendConnFailWithMsg("消息发送失败", conn)
			continue
		}

		item := ChatResponse{
			ChatListResponse: ChatListResponse{
				ChatModel:        model,
				SendUserNickname: user.Nickname,
				SendUserAvatar:   user.Avatar,
				RevUserNickname:  revUser.Nickname,
				RevUserAvatar:    revUser.Avatar,
			},
		}
		// 发给对方的
		res.SendWsMsg(OnlineMap, req.RevUserID, item)

		// 发给自己的
		item.IsMe = true
		res.SendConnOkWithData(item, conn)
	}
	defer conn.Close()

	addrMap2, ok2 := OnlineMap[userID]
	if ok2 {
		_, ok3 := addrMap2[addr]
		if ok3 {
			delete(OnlineMap[userID], addr)
		}
		if len(addrMap2) == 0 {
			delete(OnlineMap, userID)
		}
	}

	fmt.Println("离开", OnlineMap)
}
