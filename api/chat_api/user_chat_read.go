package chat_api

import (
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/gin-gonic/gin"
)

func (ChatApi) ChatReadView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	var chat models.ChatModel
	err := global.DB.Take(&chat, cr.ID).Error
	if err != nil {
		res.FailWithMsg("消息不存在", c)
		return
	}
	claims := jwts.GetClaims(c)
	var chatAc models.UserChatActionModel
	err = global.DB.Take(&chatAc, "user_id = ? and chat_id = ?", claims.UserID, cr.ID).Error
	if err != nil {
		global.DB.Create(&models.UserChatActionModel{
			UserID: claims.UserID,
			ChatID: cr.ID,
			IsRead: true,
		})
		res.OkWithMsg("消息读取成功", c)
		return
	}
	if chatAc.IsDelete {
		res.FailWithMsg("消息被删除", c)
		return
	}
	global.DB.Model(&chatAc).Update("is_read", true)
	res.OkWithMsg("消息读取成功", c)
}
