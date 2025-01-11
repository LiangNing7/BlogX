package router

import (
	"github.com/LiangNing7/BlogX/api"
	"github.com/gin-gonic/gin"
)

func SiteRouter(r *gin.RouterGroup) {
	app := api.App.SiteApi
	r.GET("site", app.SiteInfoView)
	r.PUT("site", app.SiteUpdateView)
}
