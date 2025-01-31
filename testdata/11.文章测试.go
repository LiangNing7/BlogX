package main

import (
	"fmt"

	"github.com/LiangNing7/BlogX/core"
	"github.com/LiangNing7/BlogX/flags"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.DB = core.InitDB()
	// err := global.DB.Create(&models.ArticleModel{
	// 	Title:   "嘻嘻嘻",
	// 	TagList: ctype.List{"python", "go"},
	// }).Error
	// fmt.Println(err)
	var list1 []models.ArticleModel
	global.DB.Find(&list1)
	fmt.Println(list1)
}