package log_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/LiangNing7/BlogX/core"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/gin-gonic/gin"
	e "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ActionLog struct {
	c                  *gin.Context
	level              enum.LogLevelType
	title              string
	requestBody        []byte
	responseBody       []byte
	log                *models.LogModel
	showRequestHeader  bool
	showRequest        bool
	showResponse       bool
	showResponseHeader bool
	itemList           []string
	responseHeader     http.Header
	isMiddleware       bool
}

func (ac *ActionLog) ShowRequest() {
	ac.showRequest = true
}
func (ac *ActionLog) ShowResponse() {
	ac.showResponse = true
}
func (ac *ActionLog) SetTitle(title string) {
	ac.title = title
}
func (ac *ActionLog) SetLevel(level enum.LogLevelType) {
	ac.level = level
}
func (ac *ActionLog) SetLink(label string, href string) {
	ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_item link\"><div class=\"log_item_label\">%s</div><div class=\"log_item_content\"><a href=\"%s\" target=\"_blank\">%s</a></div></div> ",
		label,
		href, href))
}
func (ac *ActionLog) SetImage(src string) {
	ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_image\"><img src=\"%s\" alt=\"\"></div>", src))
}
func (ac *ActionLog) ShowRequestHeader() {
	ac.showRequestHeader = true
}
func (ac *ActionLog) ShowResponseHeader() {
	ac.showResponseHeader = true
}
func (ac *ActionLog) setItem(label string, value any, logLevelType enum.LogLevelType) {
	var v string
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice:
		byteData, _ := json.Marshal(value)
		v = string(byteData)
	default:
		v = fmt.Sprintf("%v", value)
	}
	ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_item %s\"><div class=\"log_item_label\">%s</div><div class=\"log_item_content\">%s</div></div>",
		logLevelType,
		label, v))
}
func (ac *ActionLog) SetItem(label string, value any) {
	ac.setItem(label, value, enum.LogInfoLevel)
}
func (ac *ActionLog) SetItemInfo(label string, value any) {
	ac.setItem(label, value, enum.LogInfoLevel)
}
func (ac *ActionLog) SetItemWarn(label string, value any) {
	ac.setItem(label, value, enum.LogWarnLevel)
}
func (ac *ActionLog) SetItemError(label string, value any) {
	ac.setItem(label, value, enum.LogErrLevel)
}

func (ac *ActionLog) SetError(label string, err error) {
	msg := e.WithStack(err)
	logrus.Errorf("%s %s", label, err.Error())
	ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_error\"><div class=\"line\"><div class=\"label\">%s</div><div class=\"value\">%s</div><div class=\"type\">%T</div></div><div class=\"stack\">%+v</div></div>",
		label, err, err, msg))
}

func (ac *ActionLog) SetRequest(c *gin.Context) {
	byteData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Errorf(err.Error())
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(byteData))
	ac.requestBody = byteData
}
func (ac *ActionLog) SetResponse(data []byte) {
	ac.responseBody = data
}
func (ac *ActionLog) SetResponseHeader(header http.Header) {
	ac.responseHeader = header
}

func (ac *ActionLog) MiddlewareSave() {
	if ac.log == nil {
		// 创建
		ac.isMiddleware = true
		ac.Save()
		return
	}
	// 更新
	// 响应头
	if ac.showResponseHeader {
		byteData, _ := json.Marshal(ac.responseHeader)
		ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_response_header\"><pre class=\"log_json_body\">%s</pre></div>", string(byteData)))
	}
	// 设置响应
	if ac.showResponse {
		ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_response\"><pre class=\"log_json_body\">%s</pre></div>", string(ac.responseBody)))
	}
	ac.Save()
}

func (ac *ActionLog) Save() (id uint) {
	if ac.log != nil {
		newContent := strings.Join(ac.itemList, "\n")
		content := ac.log.Content + "\n" + newContent

		// 之前已经save过了，那就是更新
		global.DB.Model(ac.log).Updates(map[string]any{
			"content": content,
		})
		ac.itemList = []string{}
		return ac.log.ID
	}
	var newItemList []string
	// 请求头
	if ac.showRequestHeader {
		byteData, _ := json.Marshal(ac.c.Request.Header)
		newItemList = append(newItemList, fmt.Sprintf("<div class=\"log_request_header\"><pre class=\"log_json_body\">%s</pre></div>", string(byteData)))
	}
	// 设置请求
	if ac.showRequest {
		newItemList = append(newItemList, fmt.Sprintf("<div class=\"log_request\"><div class=\"log_request_head\"><span class=\"log_request_method %s\">%s</span><span class=\"log_request_path\">%s</span></div><div class=\"log_request_body\"><pre class=\"log_json_body\">%s</pre></div></div>",
			strings.ToLower(ac.c.Request.Method),
			ac.c.Request.Method,
			ac.c.Request.URL.String(),
			string(ac.requestBody),
		))
	}
	// 中间的一些content
	newItemList = append(newItemList, ac.itemList...)

	if ac.isMiddleware {
		// 响应头
		if ac.showResponseHeader {
			byteData, _ := json.Marshal(ac.responseHeader)
			newItemList = append(newItemList, fmt.Sprintf("<div class=\"log_response_header\"><pre class=\"log_json_body\">%s</pre></div>", string(byteData)))
		}
		// 设置响应
		if ac.showResponse {
			newItemList = append(newItemList, fmt.Sprintf("<div class=\"log_response\"><pre class=\"log_json_body\">%s</pre></div>", string(ac.responseBody)))
		}
	}

	ip := ac.c.ClientIP()
	addr := core.GetIpAddr(ip)
	userID := uint(1)
	log := models.LogModel{
		LogType: enum.ActionLogType,
		Title:   ac.title,
		Content: strings.Join(newItemList, "\n"),
		Level:   ac.level,
		UserID:  userID,
		IP:      ip,
		Addr:    addr,
	}
	err := global.DB.Create(&log).Error
	if err != nil {
		logrus.Errorf("日志创建失败 %s", err)
		return
	}
	ac.log = &log
	ac.itemList = []string{}
	return log.ID
}
func NewActionLogByGin(c *gin.Context) *ActionLog {
	return &ActionLog{
		c: c,
	}
}
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
