package article_api

import (
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/gin-gonic/gin"
)

type ArticleExamineRequest struct {
	ArticleID uint               `json:"articleID" binding:"required"`
	Status    enum.ArticleStatus `json:"status" binding:"required,oneof=3 4"`
	Msg       string             `json:"msg"` // 为4的时候，传递进来
}

func (ArticleApi) ArticleExamineView(c *gin.Context) {
	cr := middleware.GetBind[ArticleExamineRequest](c)
	var article models.ArticleModel
	err := global.DB.Take(&article, cr.ArticleID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}
	global.DB.Model(&article).Update("status", cr.Status)
	// TODO: 给文章的发布人发一个系统通知
	res.OkWithMsg("审核成功", c)
}
