package article_api

import (
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/ctype"
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/LiangNing7/BlogX/utils"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/gin-gonic/gin"
)

func (ArticleApi) ArticleTagOptionsView(c *gin.Context) {
	claims := jwts.GetClaims(c)
	var articleList []models.ArticleModel
	global.DB.Find(&articleList, "user_id = ? and status = ?", claims.UserID, enum.ArticleStatusPublished)
	var tagList ctype.List
	for _, model := range articleList {
		tagList = append(tagList, model.TagList...)
	}
	tagList = utils.Unique(tagList)
	var list = make([]models.OptionsResponse[string], 0)
	for _, s := range tagList {
		list = append(list, models.OptionsResponse[string]{
			Label: s,
			Value: s,
		})
	}
	res.OkWithData(list, c)
}
