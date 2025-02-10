package main

import (
	"context"
	"fmt"

	"github.com/LiangNing7/BlogX/core"
	"github.com/LiangNing7/BlogX/flags"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/olivere/elastic/v7"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.ESClient = core.EsConnect()

	query := elastic.NewBoolQuery()

	query.Must(elastic.NewMatchQuery("title", "python"))

	highlight := elastic.NewHighlight()
	highlight.Field("title")

	res, err := global.ESClient.
		Search(models.ArticleModel{}.Index()).
		Query(query).Highlight(highlight).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	count := res.Hits.TotalHits.Value // 总数
	fmt.Println(count)
	for _, hit := range res.Hits.Hits {
		fmt.Println(string(hit.Source))
		fmt.Println(hit.Highlight["title"])
	}
}