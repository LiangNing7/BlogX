package log_service

import (
	"bytes"
	"io"

	"github.com/LiangNing7/BlogX/core"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ActionLog struct {
	c            *gin.Context
	level        enum.LogLevelType
	title        string
	requestBody  []byte
	responseBody []byte
	log          *models.LogModel
}

func (ac *ActionLog) SetTitle(title string) {
	ac.title = title
}

func (ac *ActionLog) SetLevel(level enum.LogLevelType) {
	ac.level = level
}

func (ac *ActionLog) SetRequest(c *gin.Context) {
	byteData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Errorf(err.Error())
	}
	// 将读取到的请求体重新赋值给 c.Request.Body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(byteData))
	ac.requestBody = byteData
}

func (ac *ActionLog) SetResponse(data []byte) {
	ac.responseBody = data
}

func (ac *ActionLog) Save() {

	if ac.log != nil {
		// 之前已经 Save 过了，直接更新即可
		global.DB.Model(ac.log).Updates(map[string]any{
			"title": "更新",
		})
		return
	}

	ip := ac.c.ClientIP()
	addr := core.GetIpAddr(ip)
	userID := uint(1)

	log := models.LogModel{
		LogType: enum.ActionLogType,
		Title:   ac.title,
		Content: "",
		Level:   ac.level,
		UserID:  userID,
		IP:      ip,
		Addr:    addr,
	}
	err := global.DB.Create(&log).Error
	if err != nil {
		logrus.Errorf("日志创建失败 %s\n", err)
		return
	}
	ac.log = &log
}

func NewActionLogByGin(c *gin.Context) *ActionLog {
	return &ActionLog{
		c: c,
	}
}

// GetLog 防止在中间件获取一个 log，在视图中获取一个 log，
// 这里直接获取在中间件绑定在 c.Context 中的 log,没有的话再创建一个新的
func GetLog(c *gin.Context) *ActionLog {
	_log, ok := c.Get("log")
	if !ok {
		return NewActionLogByGin(c)
	}
	log, ok := _log.(*ActionLog)
	if !ok {
		return NewActionLogByGin(c)
	}
	return log
}
