package search_api

import (
	"context"
	"encoding/json"

	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/service/text_service"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type TextSearchRequest struct {
	common.PageInfo
}
type TextSearchResponse struct {
	ArticleID uint   `json:"articleID"`
	Head      string `json:"head"`
	Body      string `json:"body"`
}

func (SearchApi) TextSearchView(c *gin.Context) {
	var cr = middleware.GetBind[TextSearchRequest](c)
	if global.ESClient == nil {
		// 服务降级，用户可能没有配置es
		_list, count, _ := common.ListQuery(models.TextModel{}, common.Options{
			PageInfo: cr.PageInfo,
			Likes:    []string{"head", "body"},
		})
		var list = make([]TextSearchResponse, 0)
		for _, model := range _list {
			list = append(list, TextSearchResponse{
				ArticleID: model.ArticleID,
				Head:      model.Head,
				Body:      model.Body,
			})
		}
		res.OkWithList(list, count, c)
		return
	}
	query := elastic.NewBoolQuery()
	if cr.Key != "" {
		query.Should(
			elastic.NewMatchQuery("head", cr.Key),
			elastic.NewMatchQuery("body", cr.Key),
		)
	}
	highlight := elastic.NewHighlight()
	highlight.Field("head")
	highlight.Field("body")
	result, err := global.ESClient.
		Search(models.TextModel{}.Index()).
		Query(query).
		Highlight(highlight).
		From(cr.GetOffset()).
		Size(cr.GetLimit()).
		Do(context.Background())
	if err != nil {
		source, _ := query.Source()
		byteData, _ := json.Marshal(source)
		logrus.Errorf("查询失败 %s \n %s", err, string(byteData))
		res.FailWithMsg("查询失败", c)
		return
	}
	count := result.Hits.TotalHits.Value
	var list = make([]TextSearchResponse, 0)
	for _, hit := range result.Hits.Hits {
		var item text_service.TextModel
		err = json.Unmarshal(hit.Source, &item)
		if err != nil {
			logrus.Warnf("解析失败 %s  %s", err, string(hit.Source))
			continue
		}
		if len(hit.Highlight["head"]) > 0 {
			item.Head = hit.Highlight["head"][0]
		}
		if len(hit.Highlight["body"]) > 0 {
			item.Body = hit.Highlight["body"][0]
		}
		list = append(list, TextSearchResponse{
			ArticleID: item.ArticleID,
			Head:      item.Head,
			Body:      item.Body,
		})
	}
	res.OkWithList(list, int(count), c)
}
