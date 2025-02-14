package article_api

import (
	"fmt"

	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/LiangNing7/BlogX/service/redis_service/redis_article"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/LiangNing7/BlogX/utils/sql"
	"github.com/gin-gonic/gin"
)

type ArticleListRequest struct {
	common.PageInfo
	Type       int8               `form:"type" binding:"required,oneof=1 2 3"` // 1 用户查别人的  2 查自己的  3 管理员查
	UserID     uint               `form:"userID"`
	CategoryID *uint              `form:"categoryID"`
	Status     enum.ArticleStatus `form:"status"`
	CollectID  int                `form:"collectID"`
}

type ArticleListResponse struct {
	models.ArticleModel
	UserTop       bool    `json:"userTop"`  // 是否是用户置顶
	AdminTop      bool    `json:"adminTop"` // 是否是管理员置顶
	CategoryTitle *string `json:"categoryTitle"`
	UserNickname  string  `json:"userNickname"`
	UserAvatar    string  `json:"userAvatar"`
}

func (ArticleApi) ArticleListView(c *gin.Context) {
	cr := middleware.GetBind[ArticleListRequest](c)

	var topArticleIDList []uint // [1 2 3] => (1,2,3)

	var orderColumnMap = map[string]bool{
		"look_count desc":    true,
		"digg_count desc":    true,
		"comment_count desc": true,
		"collect_count desc": true,
		"look_count asc":     true,
		"digg_count asc":     true,
		"comment_count asc":  true,
		"collect_count asc":  true,
	}

	switch cr.Type {
	case 1:
		// 查别人。用户id就是必填的
		if cr.UserID == 0 {
			res.FailWithMsg("用户id必填", c)
			return
		}
		if cr.Page > 2 || cr.Limit > 10 {
			res.FailWithMsg("查询更多，请登录", c)
			return
		}
		cr.Status = 0
		cr.Order = ""
		if cr.CollectID != 0 {
			// 如果传了收藏夹id，那要看看这个人
			if cr.UserID == 0 {
				res.FailWithMsg("请传入用户id", c)
				return
			}

			var userConf models.UserConfModel
			err := global.DB.Take(&userConf, "user_id = ?", cr.UserID).Error
			if err != nil {
				res.FailWithMsg("用户不存在", c)
				return
			}

			if !userConf.OpenCollect {
				res.FailWithMsg("用户未开启我的收藏", c)
				return
			}
		}
	case 2:
		// 查自己的
		claims, err := jwts.ParseTokenByGin(c)
		if err != nil {
			res.FailWithMsg("请登录", c)
			return
		}
		cr.UserID = claims.UserID
	case 3:
		// 管理员
		claims, err := jwts.ParseTokenByGin(c)
		if !(err == nil && claims.Role == enum.AdminRole) {
			res.FailWithMsg("角色错误", c)
			return
		}
	}

	query := global.DB.Where("")

	if cr.CollectID != 0 {
		var articleIDList []uint
		if cr.CollectID != -1 {
			global.DB.Model(models.UserArticleCollectModel{}).Where("collect_id = ?", cr.CollectID).Select("article_id").Scan(&articleIDList)
		} else {
			// 查这个人的所有收藏夹id
			if cr.UserID == 0 {
				res.FailWithMsg("查所有的收藏文章，需要传用户id", c)
				return
			}
			global.DB.Model(models.UserArticleCollectModel{}).Where("user_id = ?", cr.UserID).Select("article_id").Scan(&articleIDList)
		}

		query.Where("id in ?", articleIDList)
	}

	if cr.Order != "" {
		_, ok := orderColumnMap[cr.Order]
		if !ok {
			res.FailWithMsg("不支持的排序方式", c)
			return
		}
	}

	var userTopMap = map[uint]bool{}
	var adminTopMap = map[uint]bool{}
	var userTopQuery = global.DB.Where("")
	if cr.UserID != 0 {
		userTopQuery.Where("user_id = ?", cr.UserID)
	}
	var userTopArticleList []models.UserTopArticleModel
	global.DB.Preload("UserModel").Order("created_at desc").Where(userTopQuery).Find(&userTopArticleList)

	for _, i2 := range userTopArticleList {
		topArticleIDList = append(topArticleIDList, i2.ArticleID)
		if i2.UserModel.Role == enum.AdminRole {
			adminTopMap[i2.ArticleID] = true
		}
		userTopMap[i2.ArticleID] = true
	}

	var options = common.Options{
		Likes:        []string{"title"},
		PageInfo:     cr.PageInfo,
		DefaultOrder: "created_at desc",
		Where:        query,
		Preloads:     []string{"CategoryModel", "UserModel"},
	}
	if len(topArticleIDList) > 0 {
		options.DefaultOrder = fmt.Sprintf("%s, created_at desc", sql.ConvertSliceOrderSql(topArticleIDList))
	}
	_list, count, _ := common.ListQuery(models.ArticleModel{
		UserID:     cr.UserID,
		CategoryID: cr.CategoryID,
		Status:     cr.Status,
	}, options)

	var list = make([]ArticleListResponse, 0)
	collectMap := redis_article.GetAllCacheCollect()
	diggMap := redis_article.GetAllCacheDigg()
	lookMap := redis_article.GetAllCacheLook()
	commentMap := redis_article.GetAllCacheComment()

	for _, model := range _list {
		model.Content = ""
		model.DiggCount = model.DiggCount + diggMap[model.ID]
		model.CollectCount = model.CollectCount + collectMap[model.ID]
		model.LookCount = model.LookCount + lookMap[model.ID]
		model.CommentCount = model.CommentCount + commentMap[model.ID]
		data := ArticleListResponse{
			ArticleModel: model,
			UserTop:      userTopMap[model.ID],
			AdminTop:     adminTopMap[model.ID],
			UserNickname: model.UserModel.Nickname,
			UserAvatar:   model.UserModel.Avatar,
		}

		if model.CategoryModel != nil {
			data.CategoryTitle = &model.CategoryModel.Title
		}

		list = append(list, data)
	}
	res.OkWithList(list, count, c)
}
