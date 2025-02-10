package ai_api

import (
	"encoding/json"

	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/service/ai_service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ArticleAnalysisRequest struct {
	Content string `json:"content" binding:"required"`
}
type ArticleAnalysisResponse struct {
	Title    string   `json:"title"`
	Abstract string   `json:"abstract"`
	Category string   `json:"category"`
	Tag      []string `json:"tag"`
}

func (AiApi) ArticleAnalysisView(c *gin.Context) {
	cr := middleware.GetBind[ArticleAnalysisRequest](c)
	if !global.Config.Ai.Enable {
		res.FailWithMsg("站点未启用ai功能", c)
		return
	}
	msg, err := ai_service.Chat(cr.Content)
	if err != nil {
		logrus.Errorf("ai分析失败 %s %s", err, cr.Content)
		res.FailWithMsg("ai分析失败", c)
		return
	}
	var data ArticleAnalysisResponse
	err = json.Unmarshal([]byte(msg), &data)
	if err != nil {
		logrus.Errorf("ai分析失败 %s %s", err, msg)
		res.FailWithMsg("ai分析失败", c)
		return
	}
	res.OkWithData(data, c)
}
