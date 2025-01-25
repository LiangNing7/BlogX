package core

import (
	"github.com/LiangNing7/BlogX/global"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

func EsConnect() *elastic.Client {
	es := global.Config.ES
	if es.Url == "" {
		return nil
	}
	client, err := elastic.NewClient(
		elastic.SetURL(es.Url),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(es.Username, es.Password),
	)
	if err != nil {
		logrus.Panicf("es连接失败 %s", err)
		return nil
	}
	logrus.Infof("es连接成功")
	return client
}
