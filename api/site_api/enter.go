package site_api

import (
	"fmt"
	"io"

	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/LiangNing7/BlogX/service/log_service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SiteApi struct {
}

func (SiteApi) SiteInfoView(c *gin.Context) {
	log_service.NewLoginSuccess(c, enum.UserPwdLoginType)
	log_service.NewLoginFail(c, enum.UserPwdLoginType, "用户不存在", "liangning", "1234")
	c.JSON(200, gin.H{"code": 0, "msg": "站点信息"})
	return
}

type SiteUpdateRequest struct {
	Name string `json:"name"`
}

func (SiteApi) SiteUpdateView(c *gin.Context) {
	log := log_service.NewActionLogByGin(c)

	// 通过绑定获取一次
	var cr SiteUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		logrus.Errorf(err.Error())
	}
	
	// 读原始 body
	byteData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Errorf(err.Error())
	}
	fmt.Println("body: ", string(byteData))

	log.Save()
	c.JSON(200, gin.H{"code": 0, "msg": "站点信息"})
	return
}
