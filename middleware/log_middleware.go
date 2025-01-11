package middleware

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	byteData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Errorf(err.Error())
	}
	fmt.Println("body: ", string(byteData))

	// 将读取到的请求体重新赋值给 c.Request.Body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(byteData))
	res := &ResponseWriter{
		ResponseWriter: c.Writer,
	}
	c.Writer = res
	c.Next()
	// 响应中间件
	fmt.Println("response: ", string(res.Body))
}
