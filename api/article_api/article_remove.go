package article_api

import (
	"fmt"

	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/gin-gonic/gin"
)

func (ArticleApi) ArticleRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.RemoveRequest](c)
	var list []models.ArticleModel
	global.DB.Find(&list, "id in ?", cr.IDList)
	if len(list) > 0 {
		err := global.DB.Delete(&list).Error
		if err != nil {
			res.FailWithMsg("删除失败", c)
			return
		}
	}
	res.OkWithMsg(fmt.Sprintf("删除成功 成功删除%d条", len(list)), c)
}
