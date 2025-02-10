package main

import (
	"fmt"

	"github.com/LiangNing7/BlogX/core"
	"github.com/LiangNing7/BlogX/flags"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/service/redis_service/redis_article"
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

	collectMap := redis_article.GetAllCacheCollect()
	diggMap := redis_article.GetAllCacheDigg()
	lookMap := redis_article.GetAllCacheLook()
	commentMap := redis_article.GetAllCacheComment()
	fmt.Printf("collectMap:%v\ndiggMap:%v\nlookMap:%v\ncommentMap:%v\n", collectMap, diggMap, lookMap, commentMap)
}
