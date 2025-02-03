package site_msg_api

import (
	"fmt"

	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum/message_type_enum"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/gin-gonic/gin"
)

type SiteMsgRemoveRequest struct {
	ID uint `json:"id"`
	T  int8 `json:"t"` // 一键已读的类型
}

func (SiteMsgApi) SiteMsgRemoveView(c *gin.Context) {
	cr := middleware.GetBind[SiteMsgRemoveRequest](c)
	claims := jwts.GetClaims(c)
	if cr.ID != 0 {
		// 找这个消息是不是当前用户的
		var msg models.MessageModel
		err := global.DB.Take(&msg, "id = ? and rev_user_id = ?", cr.ID, claims.UserID).Error
		if err != nil {
			res.FailWithMsg("消息不存在", c)
			return
		}
		global.DB.Delete(&msg)
		res.OkWithMsg("消息删除成功", c)
		return
	}
	var typeList []message_type_enum.Type
	switch cr.T {
	case 1:
		typeList = append(typeList, message_type_enum.CommentType, message_type_enum.ApplyType)
	case 2:
		typeList = append(typeList, message_type_enum.DiggArticleType, message_type_enum.DiggCommentType, message_type_enum.CollectArticleType)
	case 3:
		typeList = append(typeList, message_type_enum.SystemType)
	}
	var msgList []models.MessageModel
	global.DB.Find(&msgList, "rev_user_id = ? and type in ? ", claims.UserID, typeList)
	if len(msgList) > 0 {
		global.DB.Delete(&msgList)
	}
	res.OkWithMsg(fmt.Sprintf("批量删除%d条消息成功", len(msgList)), c)
}
