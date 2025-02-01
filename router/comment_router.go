package router

import (
	"github.com/LiangNing7/BlogX/api"
	"github.com/LiangNing7/BlogX/api/comment_api"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/gin-gonic/gin"
)

func CommentRouter(r *gin.RouterGroup) {
	app := api.App.CommentApi
	r.POST("comment", middleware.AuthMiddleware, middleware.BindJsonMiddleware[comment_api.CommentCreateRequest], app.CommentCreateView)
	r.GET("comment/tree/:id", middleware.BindUriMiddleware[models.IDRequest], app.CommentTreeView)
	r.GET("comment", middleware.AuthMiddleware, middleware.BindQueryMiddleware[comment_api.CommentListRequest], app.CommentListView)
}
