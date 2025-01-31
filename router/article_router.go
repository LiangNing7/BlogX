package router

import (
	"github.com/LiangNing7/BlogX/api"
	"github.com/LiangNing7/BlogX/api/article_api"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/gin-gonic/gin"
)

func ArticleRouter(r *gin.RouterGroup) {
	app := api.App.ArticleApi
	r.POST("article", middleware.AuthMiddleware, middleware.BindJsonMiddleware[article_api.ArticleCreateRequest], app.ArticleCreateView)
	r.GET("article", middleware.BindQueryMiddleware[article_api.ArticleListRequest], app.ArticleListView)
}
