package log_service

import (
	"github.com/LiangNing7/BlogX/core"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ActionLog struct {
	c     *gin.Context
	level enum.LogLevelType
	title string
}

func (ac *ActionLog) SetTitle(title string) {
	ac.title = title
}

func (ac *ActionLog) SetLevel(level enum.LogLevelType) {
	ac.level = level
}

func (ac *ActionLog) Save() {
	ip := ac.c.ClientIP()
	addr := core.GetIpAddr(ip)
	userID := uint(1)
	err := global.DB.Create(&models.LogModel{
		LogType: enum.ActionLogType,
		Title:   ac.title,
		Content: "",
		Level:   ac.level,
		UserID:  userID,
		IP:      ip,
		Addr:    addr,
	}).Error
	if err != nil {
		logrus.Errorf("日志创建失败 %s\n", err)
		return
	}
}

func NewActionLogByGin(c *gin.Context) *ActionLog {
	return &ActionLog{
		c: c,
	}
}
