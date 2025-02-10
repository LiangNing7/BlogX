package main

import (
	"fmt"
	"os"

	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/service/text_service"
)

func main() {
	byteData, _ := os.ReadFile("text.md")
	list := text_service.MdContentTransformation(models.ArticleModel{
		Model:   models.Model{ID: 1},
		Title:   "xxx",
		Content: string(byteData),
	})
	fmt.Println(list)
}
