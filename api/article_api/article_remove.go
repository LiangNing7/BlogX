package article_api

import (
	"fmt"

	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/service/message_service"
	"github.com/gin-gonic/gin"
)

func (ArticleApi) ArticleRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.RemoveRequest](c)
	var list []models.ArticleModel
	global.DB.Find(&list, "id in ?", cr.IDList)
	if len(list) > 0 {
		for _, model := range list {
			message_service.InsertSystemMessage(model.UserID, "管理员删除了你的文章", fmt.Sprintf("%s 内容不符合社区规范", model.Title), "", "")
		}
		err := global.DB.Delete(&list).Error
		if err != nil {
			res.FailWithMsg("删除失败", c)
			return
		}
	}
	res.OkWithMsg(fmt.Sprintf("删除成功 成功删除%d条", len(list)), c)
}
