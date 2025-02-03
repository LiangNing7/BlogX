package site_msg_api

import (
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/LiangNing7/BlogX/utils/mps"
	"github.com/gin-gonic/gin"
)

func (SiteMsgApi) UserSiteMessageConfView(c *gin.Context) {
	claims := jwts.GetClaims(c)
	var userMsgConf models.UserMessageConfModel
	err := global.DB.Take(&userMsgConf, "user_id = ?", claims.UserID).Error
	if err != nil {
		res.FailWithMsg("用户消息配置不存在", c)
		return
	}
	res.OkWithData(userMsgConf, c)
}

type UserMessageConfUpdateRequest struct {
	OpenCommentMessage *bool `json:"openCommentMessage" u:"open_comment_message"` // 开启回复和评论
	OpenDiggMessage    *bool `json:"openDiggMessage" u:"open_digg_message"`       // 开启赞和收藏
	OpenPrivateChat    *bool `json:"openPrivateChat" u:"open_private_chat"`       // 是否开启私聊
}

func (SiteMsgApi) UserSiteMessageConfUpdateView(c *gin.Context) {
	var cr = middleware.GetBind[UserMessageConfUpdateRequest](c)
	claims := jwts.GetClaims(c)
	var userMsgConf models.UserMessageConfModel
	err := global.DB.Take(&userMsgConf, "user_id = ?", claims.UserID).Error
	if err != nil {
		res.FailWithMsg("用户消息配置不存在", c)
		return
	}
	mp := mps.StructToMap(cr, "u")
	global.DB.Model(&userMsgConf).Updates(mp)
	res.OkWithMsg("用户消息配置更新成功", c)
}
