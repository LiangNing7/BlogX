package comment_api

import (
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/LiangNing7/BlogX/service/comment_service"
	"github.com/gin-gonic/gin"
)

func (CommentApi) CommentTreeView(c *gin.Context) {
	var cr = middleware.GetBind[models.IDRequest](c)
	var article models.ArticleModel
	err := global.DB.Take(&article, "status = ? and id = ?", enum.ArticleStatusPublished, cr.ID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}
	// 把根评论查出来
	var commentList []models.CommentModel
	global.DB.Find(&commentList, "article_id = ? and parent_id is null", cr.ID)
	var list = make([]comment_service.CommentResponse, 0)
	for _, model := range commentList {
		response := comment_service.GetCommentTreeV4(model.ID)
		list = append(list, *response)
	}
	res.OkWithList(list, len(list), c)
}
