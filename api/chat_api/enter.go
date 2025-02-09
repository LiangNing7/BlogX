package chat_api

import (
	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/gin-gonic/gin"
)

type ChatApi struct {
}
type ChatListRequest struct {
	common.PageInfo
	SendUserID uint `form:"sendUserID"`
	RevUserID  uint `form:"revUserID"  binding:"required"`
	Type       int8 `form:"type" binding:"required,oneof=1 2"`
}
type ChatListResponse struct {
	models.ChatModel
	SendUserNickname string `json:"sendUserNickname"`
	SendUserAvatar   string `json:"sendUserAvatar"`
	RevUserNickname  string `json:"revUserNickname"`
	RevUserAvatar    string `json:"revUserAvatar"`
	IsMe             bool   `json:"isMe"`   // 是我发的
	IsRead           bool   `json:"isRead"` // 消息是否已读
}

func (ChatApi) ChatListView(c *gin.Context) {
	cr := middleware.GetBind[ChatListRequest](c)
	claims := jwts.GetClaims(c)
	var deletedIDList []uint
	var userChatActionList []models.UserChatActionModel
	var chatReadMap = map[uint]bool{}
	global.DB.Find(&userChatActionList, "user_id = ? and (is_delete = ? or is_delete is null)", cr.RevUserID, 0)
	for _, model := range userChatActionList {
		chatReadMap[model.ChatID] = true
	}

	switch cr.Type {
	case 1: // 前台用户调的
		cr.SendUserID = claims.UserID
		// 找我删除的消息
		global.DB.Model(models.UserChatActionModel{}).
			Where("user_id = ? and is_delete = ?", claims.UserID, true).
			Select("chat_id").Scan(&deletedIDList)
	case 2:
		if claims.Role != enum.AdminRole {
			res.FailWithMsg("权限错误", c)
			return
		}
		if cr.SendUserID == 0 {
			res.FailWithMsg("sendUserID必填", c)
			return
		}
	}
	query := global.DB.Where("(send_user_id = ? and rev_user_id = ?) or(send_user_id = ? and rev_user_id = ?) ",
		cr.SendUserID, cr.RevUserID, cr.RevUserID, cr.SendUserID,
	)

	if len(deletedIDList) > 0 {
		query.Where("id not in ?", deletedIDList)
	}

	cr.Order = "created_at desc"
	_list, count, _ := common.ListQuery(models.ChatModel{}, common.Options{
		PageInfo: cr.PageInfo,
		Preloads: []string{"SendUserModel", "RevUserModel"},
		Where:    query,
	})
	var list = make([]ChatListResponse, 0)
	for _, model := range _list {
		item := ChatListResponse{
			ChatModel:        model,
			SendUserNickname: model.SendUserModel.Nickname,
			SendUserAvatar:   model.SendUserModel.Avatar,
			RevUserNickname:  model.RevUserModel.Nickname,
			RevUserAvatar:    model.RevUserModel.Nickname,
			IsRead:           chatReadMap[model.ID],
		}
		if model.SendUserID == claims.UserID {
			item.IsMe = true
		}
		list = append(list, item)
	}
	res.OkWithList(list, count, c)
}
