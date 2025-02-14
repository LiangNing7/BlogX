package search_api

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type TagAggResponse struct {
	Tag          string `json:"tag"`
	ArticleCount int    `json:"articleCount"`
}

func (SearchApi) TagAggView(c *gin.Context) {
	var cr = middleware.GetBind[common.PageInfo](c)
	var list = make([]TagAggResponse, 0)
	if global.ESClient == nil {
		var articleList []models.ArticleModel
		global.DB.Find(&articleList, "tag_list <> ''")
		var tagMap = map[string]int{}
		for _, model := range articleList {
			for _, tag := range model.TagList {
				count, ok := tagMap[tag]
				if !ok {
					tagMap[tag] = 1
					continue
				}
				tagMap[tag] = count + 1
			}
		}
		for tag, count := range tagMap {
			list = append(list, TagAggResponse{
				Tag:          tag,
				ArticleCount: count,
			})
		}
		sort.Slice(list, func(i, j int) bool {
			return list[i].ArticleCount > list[j].ArticleCount
		})
		res.OkWithList(list, len(list), c)
		return
	}
	agg := elastic.NewTermsAggregation().Field("tag_list")
	agg.SubAggregation("page",
		elastic.NewBucketSortAggregation().
			From(cr.GetOffset()).
			Size(cr.Limit))
	query := elastic.NewBoolQuery()
	query.MustNot(elastic.NewTermQuery("tag_list", ""))
	result, err := global.ESClient.
		Search(models.ArticleModel{}.Index()).
		Query(query).
		Aggregation("tags", agg).
		Aggregation("tags1", elastic.NewCardinalityAggregation().Field("tag_list")).
		Size(0).Do(context.Background())
	if err != nil {
		logrus.Errorf("查询失败 %s", err)
		res.FailWithMsg("查询失败", c)
		return
	}
	var t AggType
	var val = result.Aggregations["tags"]
	err = json.Unmarshal(val, &t)
	if err != nil {
		logrus.Errorf("解析json失败 %s %s", err, string(val))
		res.FailWithMsg("查询失败", c)
		return
	}
	var co Agg1Type
	json.Unmarshal(result.Aggregations["tags1"], &co)
	for _, bucket := range t.Buckets {
		list = append(list, TagAggResponse{
			Tag:          bucket.Key,
			ArticleCount: bucket.DocCount,
		})
	}
	res.OkWithList(list, co.Value, c)
	return
}

type AggType struct {
	DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int `json:"sum_other_doc_count"`
	Buckets                 []struct {
		Key      string `json:"key"`
		DocCount int    `json:"doc_count"`
	} `json:"buckets"`
}
type Agg1Type struct {
	Value int `json:"value"`
}
