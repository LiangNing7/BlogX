package router

import (
	"github.com/LiangNing7/BlogX/api"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/gin-gonic/gin"
)

func BannerRouter(r *gin.RouterGroup) {
	app := api.App.BannerApi
	r.GET("banner", app.BannerListView)
	r.POST("banner", middleware.AdminMiddleware, app.BannerCreateView)
	r.PUT("banner/:id", middleware.AdminMiddleware, app.BannerUpdateView)
	r.DELETE("banner", middleware.AdminMiddleware, app.BannerRemoveView)
}
