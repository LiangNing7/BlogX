package flags

import (
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/service/es_service"
)

func EsIndex() {
	article := models.ArticleModel{}
	es_service.CreateIndexV2(article.Index(), article.Mapping())
}
