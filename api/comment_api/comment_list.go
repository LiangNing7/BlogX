package comment_api

import (
	"time"

	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/LiangNing7/BlogX/models/enum/relationship_enum"
	"github.com/LiangNing7/BlogX/service/focus_service"
	"github.com/LiangNing7/BlogX/service/redis_service/redis_comment"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/gin-gonic/gin"
)

type CommentListRequest struct {
	common.PageInfo
	ArticleID uint `form:"articleID"`
	UserID    uint `form:"userID"`
	Type      int8 `form:"type" binding:"required"` // 1 查我发文章的评论  2 查我发布的评论  3 管理员看所有的评论
}
type CommentListResponse struct {
	ID           uint                       `json:"id"`
	CreatedAt    time.Time                  `json:"createdAt"`
	Content      string                     `json:"content"`
	UserID       uint                       `json:"userID"`
	UserNickname string                     `json:"userNickname"`
	UserAvatar   string                     `json:"userAvatar"`
	ArticleID    uint                       `json:"articleID"`
	ArticleTitle string                     `json:"articleTitle"`
	ArticleCover string                     `json:"articleCover"`
	DiggCount    int                        `json:"diggCount"`
	Relation     relationship_enum.Relation `json:"relation,omitempty"`
	IsMe         bool                       `json:"isMe"`
}

func (CommentApi) CommentListView(c *gin.Context) {
	cr := middleware.GetBind[CommentListRequest](c)
	query := global.DB.Where("")
	claims := jwts.GetClaims(c)
	switch cr.Type {
	case 1: // 查我发文章的评论
		// 查我发了哪些文章
		var articleIDList []uint
		global.DB.Model(models.ArticleModel{}).
			Where("user_id = ? and status = ?", claims.UserID, enum.ArticleStatusPublished).
			Select("id").Scan(&articleIDList)
		query.Where("article_id in ?", articleIDList)
		cr.UserID = 0
	case 2: // 查我发布的评论
		cr.UserID = claims.UserID
	case 3:
	}
	_list, count, _ := common.ListQuery(models.CommentModel{
		ArticleID: cr.ArticleID,
		UserID:    cr.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Likes:    []string{"content"},
		Preloads: []string{"UserModel", "ArticleModel"},
		Where:    query,
	})

	var RelationMao = map[uint]relationship_enum.Relation{}
	if cr.Type == 1 {
		var userIDList []uint
		for _, model := range _list {
			userIDList = append(userIDList, model.UserID)
		}
		RelationMao = focus_service.CalcUserPatchRelationship(claims.UserID, userIDList)
	}

	var list = make([]CommentListResponse, 0)
	for _, model := range _list {
		list = append(list, CommentListResponse{
			ID:           model.ID,
			CreatedAt:    model.CreatedAt,
			Content:      model.Content,
			UserID:       model.UserID,
			UserNickname: model.UserModel.Nickname,
			UserAvatar:   model.UserModel.Avatar,
			ArticleID:    model.ArticleID,
			ArticleTitle: model.ArticleModel.Title,
			ArticleCover: model.ArticleModel.Cover,
			DiggCount:    model.DiggCount + redis_comment.GetCacheDigg(model.ID),
			Relation:     RelationMao[model.UserID],
			IsMe:         model.UserID == claims.UserID,
		})
	}
	res.OkWithList(list, count, c)
}
