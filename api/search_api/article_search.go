package search_api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type ArticleSearchRequest struct {
	common.PageInfo
	Type int8 `form:"type"` // 0 猜你喜欢  1 最新发布  2最多回复 3最多点赞 4最多收藏
}
type ArticleBaseInfo struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Abstract string `json:"abstract"`
}
type ArticleListResponse struct {
	models.ArticleModel
	AdminTop      bool    `json:"adminTop"` // 是否是管理员置顶
	CategoryTitle *string `json:"categoryTitle"`
	UserNickname  string  `json:"userNickname"`
	UserAvatar    string  `json:"userAvatar"`
}

func (SearchApi) ArticleSearchView(c *gin.Context) {
	var cr = middleware.GetBind[ArticleSearchRequest](c)
	var sortMap = map[int8]string{
		0: "_score",
		1: "created_at",
		2: "comment_count",
		3: "digg_count",
		4: "collect_count",
	}
	sortKey := sortMap[cr.Type]
	if sortKey == "" {
		res.FailWithMsg("搜索类型错误", c)
		return
	}
	query := elastic.NewBoolQuery()
	if cr.Key != "" {
		query.Should(
			elastic.NewMatchQuery("title", cr.Key),
			elastic.NewMatchQuery("abstract", cr.Key),
			elastic.NewMatchQuery("content", cr.Key),
		)
	}
	// 只能查发布的文章
	query.Must(elastic.NewTermQuery("status", 3))
	// 把管理员置顶的文章查出来
	claims, err := jwts.ParseTokenByGin(c)
	if err == nil && claims != nil {
		// 用户登录了
		// 查用户感兴趣的分类
		var userConf models.UserConfModel
		err = global.DB.Take(&userConf, "user_id = ?", claims.UserID).Error
		if err != nil {
			res.FailWithMsg("用户配置不存在", c)
			return
		}
		if len(userConf.LikeTags) > 0 {
			tagQuery := elastic.NewBoolQuery()
			for _, tag := range userConf.LikeTags {
				tagQuery.Should(elastic.NewTermQuery("tag_list", tag))
			}
			query.Must(tagQuery)
		}
	}
	highlight := elastic.NewHighlight()
	highlight.Field("title")
	highlight.Field("abstract")
	result, err := global.ESClient.
		Search(models.ArticleModel{}.Index()).
		Query(query).
		Highlight(highlight).
		From(cr.GetOffset()).
		Size(cr.GetLimit()).
		Sort(sortKey, false).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	count := result.Hits.TotalHits.Value
	var searchArticleMap = map[uint]ArticleBaseInfo{}
	var articleIDList []uint
	for _, hit := range result.Hits.Hits {
		fmt.Println(string(hit.Source))
		fmt.Println(hit.Highlight["title"])
		fmt.Println(hit.Highlight["abstract"])
		var art ArticleBaseInfo
		err = json.Unmarshal(hit.Source, &art)
		if err != nil {
			logrus.Warnf("解析失败 %s  %s", err, string(hit.Source))
			continue
		}
		if len(hit.Highlight["title"]) > 0 {
			art.Title = hit.Highlight["title"][0]
		}
		if len(hit.Highlight["abstract"]) > 0 {
			art.Abstract = hit.Highlight["abstract"][0]
		}
		searchArticleMap[art.ID] = art
		articleIDList = append(articleIDList, art.ID)
	}
	where := global.DB.Where("id in ?", articleIDList)
	_list, _, _ := common.ListQuery(models.ArticleModel{}, common.Options{
		Where:    where,
		Preloads: []string{"CategoryModel", "UserModel"},
	})
	var list = make([]ArticleListResponse, 0)
	for _, model := range _list {
		item := ArticleListResponse{
			ArticleModel: model,
			AdminTop:     true,
			UserNickname: model.UserModel.Nickname,
			UserAvatar:   model.UserModel.Avatar,
		}
		if model.CategoryModel != nil {
			item.CategoryTitle = &model.CategoryModel.Title
		}
		item.Title = searchArticleMap[model.ID].Title
		item.Abstract = searchArticleMap[model.ID].Abstract
		list = append(list, item)
	}
	res.OkWithList(list, int(count), c)
}
