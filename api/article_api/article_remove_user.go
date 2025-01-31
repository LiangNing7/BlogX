package article_api

import (
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/gin-gonic/gin"
)

func (ArticleApi) ArticleRemoveUserView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	claims := jwts.GetClaims(c)
	var model models.ArticleModel
	err := global.DB.Take(&model, "user_id = ? and id =?", claims.UserID, cr.ID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}
	err = global.DB.Delete(&model).Error
	if err != nil {
		res.FailWithMsg("删除文章失败", c)
		return
	}
	res.OkWithMsg("删除成功", c)
}
