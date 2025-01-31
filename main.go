package main

import (
	"github.com/LiangNing7/BlogX/core"
	"github.com/LiangNing7/BlogX/flags"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/router"
	"github.com/LiangNing7/BlogX/service/cron_service"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	core.InitIPDB()
	global.DB = core.InitDB()
	global.Redis = core.InitRedis()
	global.ESClient = core.EsConnect()
	flags.Run()
	core.InitMysqlES()

	// 定时任务
	cron_service.Cron()
	
	router.Run()
}
