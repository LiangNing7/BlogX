package router

import (
	"github.com/LiangNing7/BlogX/api"
	"github.com/LiangNing7/BlogX/api/article_api"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/gin-gonic/gin"
)

func ArticleRouter(r *gin.RouterGroup) {
	app := api.App.ArticleApi
	r.POST("article", middleware.AuthMiddleware, middleware.BindJsonMiddleware[article_api.ArticleCreateRequest], app.ArticleCreateView)
	r.PUT("article", middleware.AuthMiddleware, middleware.BindJsonMiddleware[article_api.ArticleUpdateRequest], app.ArticleUpdateView)
	r.GET("article", middleware.BindQueryMiddleware[article_api.ArticleListRequest], app.ArticleListView)
	r.GET("article/:id", middleware.BindUriMiddleware[models.IDRequest], app.ArticleDetailView)
	r.POST("article/examine", middleware.AdminMiddleware, middleware.BindJsonMiddleware[article_api.ArticleExamineRequest], app.ArticleExamineView)
	r.GET("article/digg/:id", middleware.AuthMiddleware, middleware.BindUriMiddleware[models.IDRequest], app.ArticleDiggView)
}
