package site_api

import (
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/LiangNing7/BlogX/service/log_service"
	"github.com/gin-gonic/gin"
)

type SiteApi struct {
}

func (SiteApi) SiteInfoView(c *gin.Context) {
	log_service.NewLoginSuccess(c, enum.UserPwdLoginType)
	log_service.NewLoginFail(c, enum.UserPwdLoginType, "用户不存在", "liangning", "1234")
	c.JSON(200, gin.H{"code": 0, "msg": "站点信息"})
	return
}
