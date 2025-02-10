package flags

import (
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/service/es_service"
	"github.com/sirupsen/logrus"
)

func EsIndex() {
	if global.ESClient == nil {
		logrus.Warnf("未开启es连接")
		return
	}
	article := models.ArticleModel{}
	es_service.CreateIndexV2(article.Index(), article.Mapping())
	text := models.TextModel{}
	es_service.CreateIndexV2(text.Index(), text.Mapping())
}
