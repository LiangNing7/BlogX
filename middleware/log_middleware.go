package middleware

import (
	"github.com/LiangNing7/BlogX/service/log_service"
	"github.com/gin-gonic/gin"
)

// ResponseWriter 自定义的 ResponseWriter，Body用来捕获响应数据
type ResponseWriter struct {
	gin.ResponseWriter
	Body []byte
}

// Write 重写 Write 方法，捕获响应内容
func (w *ResponseWriter) Write(data []byte) (int, error) {
	w.Body = append(w.Body, data...)
	return w.ResponseWriter.Write(data)
}

func LogMiddleware(c *gin.Context) {
	// 请求中间件
	// 读原始 body
	log := log_service.NewActionLogByGin(c)
	log.SetRequest(c)
	// 绑定 log
	c.Set("log", log)

	res := &ResponseWriter{
		ResponseWriter: c.Writer,
	}
	c.Writer = res
	c.Next()
	// 响应中间件
	log.SetResponse(res.Body)
	log.Save()
}
