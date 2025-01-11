package middleware

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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
	
	c.Next()
	// 响应中间件
}
