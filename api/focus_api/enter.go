package focus_api

import (
	"time"

	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/gin-gonic/gin"
)

type FocusApi struct {
}
type FocusUserRequest struct {
	FocusUserID uint `json:"focusUserID" binding:"required"`
}

// FocusUserView 登录人关注用户
func (FocusApi) FocusUserView(c *gin.Context) {
	cr := middleware.GetBind[FocusUserRequest](c)

	claims := jwts.GetClaims(c)
	if cr.FocusUserID == claims.UserID {
		res.FailWithMsg("你时刻都在关注自己", c)
		return
	}
	// 查关注的用户是否存在
	var user models.UserModel
	err := global.DB.Take(&user, cr.FocusUserID).Error
	if err != nil {
		res.FailWithMsg("关注用户不存在", c)
		return
	}

	// 查之前是否已经关注过他了
	var focus models.UserFocusModel
	err = global.DB.Take(&focus, "user_id = ? and focus_user_id = ?", claims.UserID, user.ID).Error
	if err == nil {
		res.FailWithMsg("请勿重复关注", c)
		return
	}

	// 每天关注是不是应该有个限度？
	// 每天的取关也要有个限度？

	// 关注
	global.DB.Create(&models.UserFocusModel{
		UserID:      claims.UserID,
		FocusUserID: cr.FocusUserID,
	})

	res.OkWithMsg("关注成功", c)
	return
}

type FocusUserListRequest struct {
	common.PageInfo
	FocusUserID uint `form:"focusUserID"`
}
type FocusUserListResponse struct {
	FocusUserID       uint      `json:"focusUserID"`
	FocusUserNickname string    `json:"focusUserNickname"`
	FocusUserAvatar   string    `json:"focusUserAvatar"`
	FocusUserAbstract string    `json:"focusUserAbstract"`
	CreatedAt         time.Time `json:"createdAt"`
}

func (FocusApi) FocusUserListView(c *gin.Context) {
	cr := middleware.GetBind[FocusUserListRequest](c)

	claims := jwts.GetClaims(c)

	_list, count, _ := common.ListQuery(models.UserFocusModel{
		FocusUserID: cr.FocusUserID,
		UserID:      claims.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Preloads: []string{"FocusUserModel"},
	})

	var list = make([]FocusUserListResponse, 0)
	for _, model := range _list {
		list = append(list, FocusUserListResponse{
			FocusUserID:       model.FocusUserID,
			FocusUserNickname: model.FocusUserModel.Nickname,
			FocusUserAvatar:   model.FocusUserModel.Avatar,
			FocusUserAbstract: model.FocusUserModel.Abstract,
			CreatedAt:         model.CreatedAt,
		})
	}

	res.OkWithList(list, count, c)
}
