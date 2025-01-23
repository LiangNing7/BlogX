package user_api

import (
	"time"

	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/gin-gonic/gin"
)

type UserLoginListRequest struct {
	common.PageInfo
	UserID    uint   `form:"userId"`
	Ip        string `form:"ip"`
	Addr      string `form:"addr"`
	StartTime string `form:"startTime"` // 起止时间的 年月日时分秒格式
	EndTime   string `form:"endTime"`
	Type      int8   `form:"type" binding:"required,oneof=1 2"` // 1 用户：只能查自己的  2 管理员 ：能查全部
}
type UserLoginListResponse struct {
	models.UserLoginModel
	UserNickname string `json:"userNickname,omitempty"`
	UserAvatar   string `json:"userAvatar,omitempty"`
}

func (UserApi) UserLoginListView(c *gin.Context) {
	var cr UserLoginListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	claims := jwts.GetClaims(c)
	if cr.Type == 1 {
		cr.UserID = claims.UserID
	}
	var query = global.DB.Where("")
	if cr.StartTime != "" {
		_, err = time.Parse("2006-01-02 15:04:05", cr.StartTime)
		if err != nil {
			res.FailWithMsg("开始时间格式错误", c)
			return
		}
		query.Where("created_at >= ?", cr.StartTime)
	}
	if cr.EndTime != "" {
		_, err = time.Parse("2006-01-02 15:04:05", cr.EndTime)
		if err != nil {
			res.FailWithMsg("结束时间格式错误", c)
			return
		}
		query.Where("created_at <= ?", cr.EndTime)
	}
	var preloads []string
	if cr.Type == 2 {
		preloads = []string{"UserModel"}
	}
	_list, count, _ := common.ListQuery(models.UserLoginModel{
		UserID: cr.UserID,
		IP:     cr.Ip,
		Addr:   cr.Addr,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Where:    query,
		Preloads: preloads,
	})
	var list = make([]UserLoginListResponse, 0)
	for _, model := range _list {
		list = append(list, UserLoginListResponse{
			UserLoginModel: model,
			UserNickname:   model.UserModel.Nickname,
			UserAvatar:     model.UserModel.Avatar,
		})
	}
	res.OkWithList(list, count, c)
}
