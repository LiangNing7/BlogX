package cron_service

import (
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/service/redis_service/redis_site"
)

// SyncSiteFlow 同步网站访问量
func SyncSiteFlow() {
	flow := redis_site.GetFlow()
	global.DB.Create(&models.SiteFlowModel{Count: flow})
	redis_site.ClearFlow()
}
