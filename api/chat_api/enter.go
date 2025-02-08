package chat_api

import (
	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/gin-gonic/gin"
)

type ChatApi struct {
}
type ChatListRequest struct {
	common.PageInfo
	UserID uint `form:"userID" binding:"required"` // 查我和他的聊天记录
}
type ChatListResponse struct {
	models.ChatModel
	SendUserNickname string `json:"sendUserNickname"`
	SendUserAvatar   string `json:"sendUserAvatar"`
	RevUserNickname  string `json:"revUserNickname"`
	RevUserAvatar    string `json:"revUserAvatar"`
	IsMe             bool   `json:"isMe"` // 是我发的
}

func (ChatApi) ChatListView(c *gin.Context) {
	cr := middleware.GetBind[ChatListRequest](c)
	claims := jwts.GetClaims(c)
	query := global.DB.Where("(send_user_id = ? and rev_user_id = ?) or(send_user_id = ? and rev_user_id = ?) ",
		cr.UserID, claims.UserID, claims.UserID, cr.UserID,
	)
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
		}
		if model.SendUserID == claims.UserID {
			item.IsMe = true
		}
		list = append(list, item)
	}
	res.OkWithList(list, count, c)
}
