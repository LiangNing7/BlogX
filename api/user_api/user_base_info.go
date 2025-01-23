package user_api

import (
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/gin-gonic/gin"
)

type UserBaseInfoResponse struct {
	UserID       uint   `json:"userID"`
	CodeAge      int    `json:"codeAge"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
	LookCount    int    `json:"lookCount"`
	ArticleCount int    `json:"articleCount"`
	FansCount    int    `json:"fansCount"`
	FollowCount  int    `json:"followCount"`
	Place        string `json:"place"` // ip归属地
}

func (UserApi) UserBaseInfoView(c *gin.Context) {
	var cr models.IDRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	var user models.UserModel
	err = global.DB.Take(&user, cr.ID).Error
	if err != nil {
		res.FailWithMsg("不存在的用户", c)
		return
	}
	data := UserBaseInfoResponse{
		UserID:       user.ID,
		CodeAge:      user.CodeAge(),
		Avatar:       user.Avatar,
		Nickname:     user.Nickname,
		LookCount:    1,
		ArticleCount: 1, // TODO: 把文章做完回来做
		FansCount:    1, // TODO: 把好友关系做完回来做
		FollowCount:  1,
		Place:        user.Addr,
	}
	res.OkWithData(data, c)
}
