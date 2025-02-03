package comment_api

import (
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/service/message_service"
	"github.com/LiangNing7/BlogX/service/redis_service/redis_comment"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/gin-gonic/gin"
)

func (CommentApi) CommentDiggView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	var comment models.CommentModel
	err := global.DB.Take(&comment, cr.ID).Error
	if err != nil {
		res.FailWithMsg("评论不存在", c)
		return
	}
	claims := jwts.GetClaims(c)
	// 查一下之前有没有点过
	var userDiggComment models.CommentDiggModel
	err = global.DB.Take(&userDiggComment, "user_id = ? and comment_id = ?", claims.UserID, comment.ID).Error
	if err != nil {
		// 点赞
		model := models.CommentDiggModel{
			UserID:    claims.UserID,
			CommentID: cr.ID,
		}
		err = global.DB.Create(&model).Error
		if err != nil {
			res.FailWithMsg("点赞失败", c)
			return
		}
		redis_comment.SetCacheDigg(cr.ID, 1)
		// 给这个评论的拥有人发消息
		message_service.InsertDiggCommentMessage(model)
		res.OkWithMsg("点赞成功", c)
		return
	}
	// 取消点赞
	global.DB.Model(models.CommentDiggModel{}).Delete("user_id = ? and comment_id = ?", claims.UserID, comment.ID)
	res.OkWithMsg("取消点赞成功", c)
	redis_comment.SetCacheDigg(cr.ID, -1)
	return
}
