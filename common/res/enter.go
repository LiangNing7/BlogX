package res

import (
	"github.com/LiangNing7/BlogX/utils/validate"
	"github.com/gin-gonic/gin"
)

type Code int

const (
	SuccessCode     Code = 0
	FailValidCode   Code = 1001
	FailServiceCode Code = 1002
)

func (c Code) String() string {
	switch c {
	case SuccessCode:
		return "成功"
	case FailValidCode:
		return "校验失败"
	case FailServiceCode:
		return "服务异常"
	}
	return ""
}

var empty = map[string]any{}

type Response struct {
	Code Code   `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func (r Response) Json(c *gin.Context) {
	c.JSON(200, r)
}
func Ok(data any, msg string, c *gin.Context) {
	Response{SuccessCode, data, msg}.Json(c)
}
func OkWithData(data any, c *gin.Context) {
	Response{SuccessCode, data, "成功"}.Json(c)
}
func OkWithList(list any, count int, c *gin.Context) {
	Response{SuccessCode, map[string]any{
		"list":  list,
		"count": count,
	}, "成功"}.Json(c)
}
func OkWithMsg(msg string, c *gin.Context) {
	Response{SuccessCode, empty, msg}.Json(c)
}
func FailWithMsg(msg string, c *gin.Context) {
	Response{FailValidCode, empty, msg}.Json(c)
}
func FailWithData(data any, msg string, c *gin.Context) {
	Response{FailServiceCode, data, msg}.Json(c)
}
func FailWithCode(code Code, c *gin.Context) {
	Response{code, empty, code.String()}.Json(c)
}
func FailWithError(err error, c *gin.Context) {
	data, msg := validate.ValidateError(err)
	FailWithData(data, msg, c)
}
