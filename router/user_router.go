package router

import (
	"github.com/LiangNing7/BlogX/api"
	"github.com/LiangNing7/BlogX/api/user_api"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup) {
	app := api.App.UserApi
	r.POST("user/send_email", middleware.CaptchaMiddleware, app.SendEmailView)
	r.POST("user/email", middleware.EmailVerifyMiddleware, app.RegisterEmailView)
	r.POST("user/qq", app.QQLoginView)
	r.POST("user/login", middleware.CaptchaMiddleware, middleware.BindJsonMiddleware[user_api.PwdLoginRequest], app.PwdLoginApi)
	r.GET("user/detail", middleware.AuthMiddleware, app.UserDetailView)
	r.GET("user", middleware.AdminMiddleware, middleware.BindQueryMiddleware[user_api.UserListRequest], app.UserListView)
	r.GET("user/login", middleware.AuthMiddleware, app.UserLoginListView)
	r.GET("user/base", app.UserBaseInfoView)
	r.PUT("user/password", middleware.AuthMiddleware, app.UpdatePasswordView)
	r.PUT("user/password/reset", middleware.EmailVerifyMiddleware, app.ResetPasswordView)
	r.PUT("user/email/bind", middleware.EmailVerifyMiddleware, middleware.AuthMiddleware, app.BindEmailView)
	r.PUT("user", middleware.AuthMiddleware, app.UserInfoUpdateView)
	r.PUT("user/admin", middleware.AdminMiddleware, app.AdminUserInfoUpdateView)
	r.DELETE("user/logout", middleware.AuthMiddleware, app.LogoutView)
	r.POST("user/article/top", middleware.AuthMiddleware, middleware.BindJsonMiddleware[user_api.UserArticleTopRequest], app.UserArticleTopView)
}
