package router

import (
	"github.com/LiangNing7/BlogX/api"
	"github.com/LiangNing7/BlogX/api/data_api"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/gin-gonic/gin"
)

func DataRouter(r *gin.RouterGroup) {
	app := api.App.DataApi
	r.GET("data/sum", middleware.AdminMiddleware, app.SumView)
	r.GET("data/article/year", middleware.AdminMiddleware, app.ArticleYearDataView)
	r.GET("data/growth", middleware.AdminMiddleware, middleware.BindQueryMiddleware[data_api.GrowthDataRequest], app.GrowthDataView)
}
