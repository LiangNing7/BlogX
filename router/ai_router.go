package router

import (
	"github.com/LiangNing7/BlogX/api"
	"github.com/LiangNing7/BlogX/api/ai_api"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/gin-gonic/gin"
)

func AiRouter(r *gin.RouterGroup) {
	app := api.App.AiApi
	r.POST("ai/analysis", middleware.AuthMiddleware, middleware.BindJsonMiddleware[ai_api.ArticleAnalysisRequest], app.ArticleAnalysisView)
	r.POST("ai/article", middleware.AuthMiddleware, middleware.BindJsonMiddleware[ai_api.ArticleAiRequest], app.ArticleAiView)
}
