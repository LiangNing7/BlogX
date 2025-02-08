package router

import (
	"github.com/LiangNing7/BlogX/api"
	"github.com/LiangNing7/BlogX/api/focus_api"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/gin-gonic/gin"
)

func FocusRouter(r *gin.RouterGroup) {
	app := api.App.FocusApi
	r.POST("focus", middleware.AuthMiddleware, middleware.BindJsonMiddleware[focus_api.FocusUserRequest], app.FocusUserView)
	r.GET("focus/my_focus", middleware.BindQueryMiddleware[focus_api.FocusUserListRequest], app.FocusUserListView)
	r.GET("focus/my_fans", middleware.BindQueryMiddleware[focus_api.FocusUserListRequest], app.FansUserListView)
}
