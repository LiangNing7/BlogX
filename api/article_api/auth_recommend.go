package article_api

import (
	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum/relationship_enum"
	"github.com/LiangNing7/BlogX/service/focus_service"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/gin-gonic/gin"
)

type AuthRecommendResponse struct {
	UserID       uint   `json:"userID"`
	UserNickname string `json:"userNickname"`
	UserAvatar   string `json:"userAvatar"`
	UserAbstract string `json:"userAbstract"`
}

func (ArticleApi) AuthRecommendView(c *gin.Context) {
	cr := middleware.GetBind[common.PageInfo](c)
	var count int
	var userIDList []uint
	global.DB.Model(models.ArticleModel{}).Group("user_id").Select("count(*)").Scan(&count)
	global.DB.Model(models.ArticleModel{}).Group("user_id").
		Offset(cr.GetOffset()).
		Limit(cr.GetLimit()).
		Select("user_id").Scan(&userIDList)
	claims, err := jwts.ParseTokenByGin(c)
	if err == nil && claims != nil {
		m := focus_service.CalcUserPatchRelationship(claims.UserID, userIDList)
		userIDList = []uint{}
		for u, relation := range m {
			if relation == relationship_enum.RelationStranger || relation == relationship_enum.RelationFans {
				userIDList = append(userIDList, u)
			}
		}
	}
	var userList []models.UserModel
	global.DB.Find(&userList, "id in ?", userIDList)
	var list = make([]AuthRecommendResponse, 0)
	for _, model := range userList {
		list = append(list, AuthRecommendResponse{
			UserID:       model.ID,
			UserNickname: model.Nickname,
			UserAvatar:   model.Avatar,
			UserAbstract: model.Abstract,
		})
	}
	res.OkWithList(list, count, c)
}
