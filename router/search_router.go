package router

import (
	"github.com/LiangNing7/BlogX/api"
	"github.com/LiangNing7/BlogX/api/search_api"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/gin-gonic/gin"
)

func SearchRouter(r *gin.RouterGroup) {
	app := api.App.SearchApi
	r.GET("search/article", middleware.BindQueryMiddleware[search_api.ArticleSearchRequest], app.ArticleSearchView)
	r.GET("search/text", middleware.BindQueryMiddleware[search_api.TextSearchRequest], app.TextSearchView)
}
