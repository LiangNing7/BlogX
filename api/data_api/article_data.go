package data_api

import (
	"time"

	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/gin-gonic/gin"
)

type ArticleDataResponse struct {
	GrowthRate int      `json:"growthRate"` // 增长率
	GrowthNum  int      `json:"growthNum"`  // 增长数
	CountList  []int    `json:"countList"`
	DateList   []string `json:"dateList"`
}

func (DataApi) ArticleDataView(c *gin.Context) {
	// 1 2 3 4 5 6 7
	// 1 10%
	now := time.Now()
	// 七天前的时间
	before7 := now.AddDate(0, 0, -7)
	// 查询七天内的文章
	var articleList []models.ArticleModel
	global.DB.Find(&articleList,
		"created_at >= ? and created_at <= ? and status = ?",
		before7.Format("2006-01-02")+" 00:00:00",
		now.Format("2006-01-02 15:04:05"),
		enum.ArticleStatusPublished,
	)
	var dateMap = map[string]int{}
	for _, model := range articleList {
		date := model.CreatedAt.Format("2006-01-02")
		count, ok := dateMap[date]
		if !ok {
			dateMap[date] = 1
			continue
		}
		dateMap[date] = count + 1
	}
	response := ArticleDataResponse{}
	for i := 0; i < 7; i++ {
		date := before7.AddDate(0, 0, i+1)
		dateS := date.Format("2006-01-02")
		count, _ := dateMap[dateS]
		response.CountList = append(response.CountList, count)
		response.DateList = append(response.DateList, dateS)
	}
	// 算增长，找最后一个和最后一个的前一个
	response.GrowthNum = response.CountList[6] - response.CountList[5]
	response.GrowthRate = int(float64(response.GrowthNum) / float64(response.CountList[5]) * 100)
	res.OkWithData(response, c)
}
