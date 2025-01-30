package main

import (
	"context"
	"fmt"
	"time"

	"github.com/LiangNing7/BlogX/core"
	"github.com/LiangNing7/BlogX/flags"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/olivere/elastic/v7"
)

func create() {
	timeStr := "2025-01-31 00:17:11.000"
	layout := "2006-01-02 15:04:05.000"
	t, err := time.Parse(layout, timeStr)
	var article = models.ArticleModel{
		Model: models.Model{
			ID:        1,
			CreatedAt: t,
			UpdatedAt: t,
		},
		Title:   "LiangNing's Blog",
		Content: "这是内容",
		UserID:  1,
		Status:  1,
	}
	indexResponse, err := global.ESClient.Index().Index(article.Index()).BodyJson(article).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", indexResponse)
}
func list() {
	limit := 2
	page := 1
	from := (page - 1) * limit
	query := elastic.NewBoolQuery()
	res, err := global.ESClient.Search(models.ArticleModel{}.Index()).Query(query).From(from).Size(limit).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	count := res.Hits.TotalHits.Value // 总数
	fmt.Println(count)
	for _, hit := range res.Hits.Hits {
		fmt.Println(string(hit.Source))
	}
}
func DocDelete() {
	deleteResponse, err := global.ESClient.Delete().
		Index(models.ArticleModel{}.Index()).Id("V9kUnpQBujP7135xix7G").Refresh("true").Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(deleteResponse)
}
func update() {
	updateResponse, err := global.ESClient.Update().Index(models.ArticleModel{}.Index()).Refresh("true").
		Id("T14oB5MBjX_XsWH00KRt").
		Doc(map[string]any{
			"content": "枫枫1111",
		}).Do(context.Background())
	fmt.Println(updateResponse, err)
}
func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.ESClient = core.EsConnect()
	create()
	// list()
	// DocDelete()
	// update()
}
