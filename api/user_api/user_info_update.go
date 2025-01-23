package user_api

import (
	"fmt"
	"time"

	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/LiangNing7/BlogX/utils/mps"
	"github.com/gin-gonic/gin"
)

type UserInfoUpdateRequest struct {
	Username    *string   `json:"username" s-u:"username"`
	Nickname    *string   `json:"nickname" s-u:"nickname"`
	Avatar      *string   `json:"avatar" s-u:"avatar"`
	Abstract    *string   `json:"abstract" s-u:"abstract"`
	LikeTags    *[]string `json:"likeTags" s-u-c:"like_tags"`
	OpenCollect *bool     `json:"openCollect" s-u-c:"open_collect"`  // 公开我的收藏
	OpenFollow  *bool     `json:"openFollow" s-u-c:"open_follow"`    // 公开我的关注
	OpenFans    *bool     `json:"openFans" s-u-c:"open_fans"`        // 公开我的粉丝
	HomeStyleID *uint     `json:"homeStyleID" s-u-c:"home_style_id"` // 主页样式的id
}

func (UserApi) UserInfoUpdateView(c *gin.Context) {
	var cr UserInfoUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	userMap := mps.StructToMap(cr, "s-u")
	userConfMap := mps.StructToMap(cr, "s-u-c")
	fmt.Println("userMap", userMap)
	fmt.Println("userConfMap", userConfMap)
	claims := jwts.GetClaims(c)
	if len(userMap) > 0 {
		var userModel models.UserModel
		err = global.DB.Preload("UserConfModel").Take(&userModel, claims.UserID).Error
		if err != nil {
			res.FailWithMsg("用户不存在", c)
			return
		}
		// 判断
		if cr.Username != nil {
			var userCount int64
			global.DB.Debug().Model(models.UserModel{}).
				Where("username = ? and id <> ?", *cr.Username, claims.UserID).
				Count(&userCount)
			fmt.Println(*cr.Username, userCount)
			if userCount > 0 {
				res.FailWithMsg("该用户名被使用", c)
				return
			}
			if *cr.Username != userModel.Username {
				// 如果和我的用户名是一样的
				var uud = userModel.UserConfModel.UpdateUsernameDate
				if uud != nil {
					if time.Now().Sub(*uud).Hours() < 720 {
						res.FailWithMsg("用户名30天内只能修改一次", c)
						return
					}
				}
				userConfMap["update_username_date"] = time.Now()
			}
		}
		if cr.Nickname != nil || cr.Avatar != nil {
			if userModel.RegisterSource == enum.RegisterQQSourceType {
				res.FailWithMsg("QQ注册的用户不能修改昵称和头像", c)
				return
			}
		}
		err = global.DB.Model(&userModel).Updates(userMap).Error
		if err != nil {
			res.FailWithMsg("用户信息修改失败", c)
			return
		}
	}
	if len(userConfMap) > 0 {
		var userConfModel models.UserConfModel
		err = global.DB.Take(&userConfModel, "user_id = ?", claims.UserID).Error
		if err != nil {
			res.FailWithMsg("用户配置信息不存在", c)
			return
		}
		err = global.DB.Model(&userConfModel).Updates(userConfMap).Error
		if err != nil {
			res.FailWithMsg("用户信息修改失败", c)
			return
		}
	}
	res.OkWithMsg("用户信息修改成功", c)
}
