package user_api

import (
	"time"

	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/gin-gonic/gin"
)

type UserListRequest struct {
	common.PageInfo
}
type UserListResponse struct {
	ID            uint       `json:"id"`
	Nickname      string     `json:"nickname"`
	Username      string     `json:"username"`
	Avatar        string     `json:"avatar"`
	IP            string     `json:"ip"`
	Addr          string     `json:"addr"`
	ArticleCount  int        `json:"articleCount"`  // 发文数
	IndexCount    int        `json:"indexCount"`    // 主页访问数
	CreatedAt     time.Time  `json:"createdAt"`     // 注册时间
	LastLoginDate *time.Time `json:"lastLoginDate"` // 最后登录时间
}

func (UserApi) UserListView(c *gin.Context) {
	cr := middleware.GetBind[UserListRequest](c)
	_list, count, _ := common.ListQuery(models.UserModel{}, common.Options{
		Likes:    []string{"nickname", "username"},
		Preloads: []string{"ArticleList", "LoginList"},
		PageInfo: cr.PageInfo,
	})
	var list = make([]UserListResponse, 0)
	for _, model := range _list {
		item := UserListResponse{
			ID:           model.ID,
			Nickname:     model.Nickname,
			Username:     model.Username,
			Avatar:       model.Avatar,
			IP:           model.IP,
			Addr:         model.Addr,
			ArticleCount: len(model.ArticleList),
			IndexCount:   1000,
			CreatedAt:    model.CreatedAt,
		}
		if len(model.LoginList) > 0 {
			item.LastLoginDate = &model.LoginList[len(model.LoginList)-1].CreatedAt
		}
		list = append(list, item)
	}
	res.OkWithList(list, count, c)
}
