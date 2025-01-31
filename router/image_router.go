package router

import (
	"github.com/LiangNing7/BlogX/api"
	"github.com/LiangNing7/BlogX/api/image_api"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/gin-gonic/gin"
)

func ImageRouter(r *gin.RouterGroup) {
	app := api.App.ImageApi
	r.POST("images", middleware.AuthMiddleware, app.ImageUploadView)
	r.POST("images/qiniu", middleware.AuthMiddleware, app.QiNiuGenToken)
	r.POST("images/transfer_deposit", middleware.AuthMiddleware, middleware.BindJsonMiddleware[image_api.TransferDepositRequest], app.TransferDepositView)
	r.GET("images", middleware.AdminMiddleware, app.ImageListView)
	r.DELETE("images", middleware.AdminMiddleware, app.ImageRemoveView)
}
