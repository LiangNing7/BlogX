package main

import (
	"github.com/LiangNing7/BlogX/core"
	"github.com/LiangNing7/BlogX/flags"
	"github.com/LiangNing7/BlogX/global"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.DB = core.InitDB()

	flags.Run()
}
