package router

import (
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/gin-gonic/gin"
)

func Run() {
	gin.SetMode(global.Config.System.GinMode)
	r := gin.Default()

	// 配置静态文件路径
	r.Static("/uploads", "uploads")
	nr := r.Group("/api")
	nr.Use(middleware.LogMiddleware)
	SiteRouter(nr)
	LogRouter(nr)
	ImageRouter(nr)
	addr := global.Config.System.Addr()
	r.Run(addr)
}
