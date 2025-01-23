package user_api

import (
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/LiangNing7/BlogX/utils/pwd"
	"github.com/gin-gonic/gin"
)

type ResetPasswordRequest struct {
	Pwd string `json:"pwd" binding:"required"`
}

func (UserApi) ResetPasswordView(c *gin.Context) {
	var cr ResetPasswordRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	if !global.Config.Site.Login.EmailLogin {
		res.FailWithMsg("站点未启用邮箱注册", c)
		return
	}
	_email, _ := c.Get("email")
	email := _email.(string)
	var user models.UserModel
	err = global.DB.Take(&user, "email = ?", email).Error
	if err != nil {
		res.FailWithMsg("不存在的用户", c)
		return
	}
	if user.RegisterSource != enum.RegisterEmailSourceType {
		res.FailWithMsg("非邮箱注册用户，不能重置密码", c)
		return
	}
	hashPwd, _ := pwd.GenerateFromPassword(cr.Pwd)
	global.DB.Model(&user).Update("password", hashPwd)
	res.OkWithMsg("重置密码成功", c)
}
