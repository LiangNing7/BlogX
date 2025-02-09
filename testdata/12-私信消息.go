package main

import (
	"github.com/LiangNing7/BlogX/core"
	"github.com/LiangNing7/BlogX/flags"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/service/chat_service"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.DB = core.InitDB()
	// chat_service.ToImageChat(2, 1, "http://sqa480fov.sabkt.gdipper.com/blogx/jk雷神.jpg")
	chat_service.ToTextChat(1, 2, "你好呀!")
}
