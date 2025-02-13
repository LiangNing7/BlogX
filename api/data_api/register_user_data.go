package data_api

import (
	"time"

	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/gin-gonic/gin"
)

type RegisterUserDataResponse struct {
	GrowthRate int      `json:"growthRate"` // 增长率
	GrowthNum  int      `json:"growthNum"`  // 增长数
	CountList  []int    `json:"countList"`
	DateList   []string `json:"dateList"`
}

func (DataApi) RegisterUserDataView(c *gin.Context) {
	now := time.Now()
	before7 := now.AddDate(0, 0, -7)
	var userList []models.UserModel
	global.DB.Find(&userList,
		"created_at >= ? and created_at <= ?",
		before7.Format("2006-01-02")+" 00:00:00",
		now.Format("2006-01-02 15:04:05"),
	)
	var dateMap = map[string]int{}
	for _, model := range userList {
		date := model.CreatedAt.Format("2006-01-02")
		count, ok := dateMap[date]
		if !ok {
			dateMap[date] = 1
			continue
		}
		dateMap[date] = count + 1
	}
	response := RegisterUserDataResponse{}
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
