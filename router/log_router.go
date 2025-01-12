package router

import (
	"github.com/LiangNing7/BlogX/api"
	"github.com/gin-gonic/gin"
)

func LogRouter(r *gin.RouterGroup) {
	app := api.App.LogApi
	r.GET("logs", app.LogListView)
	r.GET("logs/:id", app.LogReadView)
}
