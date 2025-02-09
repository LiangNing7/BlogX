package chat_api

import (
	"fmt"
	"time"

	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/ctype/chat_msg"
	"github.com/LiangNing7/BlogX/models/enum/chat_msg_type"
	"github.com/LiangNing7/BlogX/models/enum/relationship_enum"
	"github.com/LiangNing7/BlogX/service/focus_service"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/LiangNing7/BlogX/utils/sql"
	"github.com/gin-gonic/gin"
)

type SessionListRequest struct {
	common.PageInfo
}

type SessionTable struct {
	SU        uint   `gorm:"column:sU"`
	RU        uint   `gorm:"column:rU"`
	MaxDate   string `gorm:"column:maxDate"`
	Count     int    `gorm:"column:c"`
	NewChatID uint   `gorm:"column:newChatID"`
}

type SessionListResponse struct {
	UserID       uint                       `json:"userID"`
	UserNickname string                     `json:"userNickname"`
	UserAvatar   string                     `json:"userAvatar"`
	Msg          chat_msg.ChatMsg           `json:"msg"`
	MsgType      chat_msg_type.MsgType      `json:"msgType"`
	NewMsgDate   time.Time                  `json:"newMsgDate"`
	Relation     relationship_enum.Relation `json:"relation"` // 好友关系
}

func (ChatApi) SessionListView(c *gin.Context) {

	cr := middleware.GetBind[SessionListRequest](c)
	claims := jwts.GetClaims(c)

	// 查我删了哪些消息
	var deletedIDList []uint
	global.DB.Model(models.UserChatActionModel{}).
		Where("user_id = ? and is_delete = ?", claims.UserID, true).
		Select("chat_id").Scan(&deletedIDList)

	query := global.DB.Where("")
	var column = fmt.Sprintf("(select id from chat_models where ((send_user_id = sU and rev_user_id = rU) or (send_user_id = rU and rev_user_id = sU))  order by created_at desc limit 1) as newChatID")
	if len(deletedIDList) > 0 {
		query.Where("id not in ?", deletedIDList)
		column = fmt.Sprintf("(select id from chat_models where ((send_user_id = sU and rev_user_id = rU) or (send_user_id = rU and rev_user_id = sU)) and id not in %s order by created_at desc limit 1) as newChatID", sql.ConvertSliceSql(deletedIDList))
	}

	var _list []SessionTable
	global.DB.Model(models.ChatModel{}).
		Select(
			"least(send_user_id, rev_user_id)    as sU",
			"greatest(send_user_id, rev_user_id) as rU",
			"max(created_at)       as maxDate",
			"count(*)         as c",
			column,
		).
		Where(query).
		Where("(send_user_id = ? or rev_user_id = ?)", claims.UserID, claims.UserID).
		Group("least(send_user_id, rev_user_id)").
		Group("greatest(send_user_id, rev_user_id)").
		Order("maxDate desc").
		Limit(cr.GetLimit()).Offset(cr.GetOffset()).Scan(&_list)

	var count int
	global.DB.Select("count(*)").Table("(?) as x",
		global.DB.
			Model(models.ChatModel{}).
			Select("count(*)").
			Where("(send_user_id = ? or rev_user_id = ?)", claims.UserID, claims.UserID).
			Group("least(send_user_id, rev_user_id)").
			Group("greatest(send_user_id, rev_user_id)"),
	).Scan(&count)

	var userIDList []uint
	var chatIDList []uint
	for _, table := range _list {
		chatIDList = append(chatIDList, table.NewChatID)
		if table.RU == claims.UserID {
			userIDList = append(userIDList, table.SU)
		}
		if table.SU == claims.UserID {
			userIDList = append(userIDList, table.RU)
		}
	}

	userMap := common.ScanMapV2(models.UserModel{}, common.ScanOption{
		Where: global.DB.Where("id in ?", userIDList),
	})
	chatMap := common.ScanMapV2(models.ChatModel{}, common.ScanOption{
		Where: global.DB.Where("id in ?", chatIDList),
	})

	relationMap := focus_service.CalcUserPatchRelationship(claims.UserID, userIDList)

	var list = make([]SessionListResponse, 0)
	for _, table := range _list {
		item := SessionListResponse{}
		if table.RU == claims.UserID {
			item.UserID = table.SU
		}
		if table.SU == claims.UserID {
			item.UserID = table.RU
		}

		item.UserNickname = userMap[item.UserID].Nickname
		item.UserAvatar = userMap[item.UserID].Avatar
		item.Msg = chatMap[table.NewChatID].Msg
		item.MsgType = chatMap[table.NewChatID].MsgType
		item.NewMsgDate = chatMap[table.NewChatID].CreatedAt
		item.Relation = relationMap[item.UserID]
		list = append(list, item)
	}

	res.OkWithList(list, count, c)
}
