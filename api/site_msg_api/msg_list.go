package site_msg_api

import (
	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum/message_type_enum"
	"github.com/LiangNing7/BlogX/models/enum/relationship_enum"
	"github.com/LiangNing7/BlogX/service/focus_service"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/gin-gonic/gin"
)

type SiteMsgListRequest struct {
	common.PageInfo
	T int8 `form:"t" binding:"required,oneof=1 2 3"` // 1评论和回复 2赞和收藏 3 系统
}

type SiteMsgListResponse struct {
	models.MessageModel
	Relation relationship_enum.Relation `json:"relation"`
}

func (SiteMsgApi) SiteMsgListView(c *gin.Context) {
	cr := middleware.GetBind[SiteMsgListRequest](c)
	var typeList []message_type_enum.Type
	switch cr.T {
	case 1:
		typeList = append(typeList, message_type_enum.CommentType, message_type_enum.ApplyType)
	case 2:
		typeList = append(typeList, message_type_enum.DiggArticleType, message_type_enum.DiggCommentType, message_type_enum.CollectArticleType)
	case 3:
		typeList = append(typeList, message_type_enum.SystemType)
	}
	claims := jwts.GetClaims(c)
	_list, count, _ := common.ListQuery(models.MessageModel{
		RevUserID: claims.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Where:    global.DB.Where("type in ?", typeList),
	})

	var userIDList []uint
	for _, model := range _list {
		if model.ActionUserID != 0 {
			userIDList = append(userIDList, model.ActionUserID)
		}
	}
	var m = map[uint]relationship_enum.Relation{}
	if len(userIDList) > 0 {
		m = focus_service.CalcUserPatchRelationship(claims.UserID, userIDList)
	}
	var list = make([]SiteMsgListResponse, 0)
	for _, model := range _list {
		list = append(list, SiteMsgListResponse{
			MessageModel: model,
			Relation:     m[model.ActionUserID],
		})
	}
	res.OkWithList(list, count, c)
}
