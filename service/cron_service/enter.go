package cron_service

import (
	"time"

	"github.com/robfig/cron/v3"
)

func Cron() {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	crontab := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))
	// 每天 4 点去同步文章数据
	crontab.AddFunc("0 0 4 * * *", SyncArticle)
	crontab.AddFunc("0 0 4 * * *", SyncComment)
	crontab.AddFunc("0 0 4 * * *", SyncSiteFlow)
	crontab.Start()
}
