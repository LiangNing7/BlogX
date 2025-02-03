package global_notification_api

import (
	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/gin-gonic/gin"
)

type GlobalNotificationApi struct {
}
type CreateRequest struct {
	Title   string `json:"title" binding:"required"`
	Icon    string `json:"icon"`
	Content string `json:"content" binding:"required"`
	Href    string `json:"href"` // 用户点击消息，然后去进行一个跳转
}

func (GlobalNotificationApi) CreateView(c *gin.Context) {
	cr := middleware.GetBind[CreateRequest](c)
	var model models.GlobalNotificationModel
	err := global.DB.Take(&model, "title = ?", cr.Title).Error
	if err == nil {
		res.FailWithMsg("全局消息名称重复", c)
		return
	}
	err = global.DB.Create(&models.GlobalNotificationModel{
		Title:   cr.Title,
		Icon:    cr.Icon,
		Content: cr.Content,
		Href:    cr.Href,
	}).Error
	if err != nil {
		res.FailWithMsg("全局消息创建失败", c)
		return
	}
	res.OkWithMsg("消息创建成功", c)
}

type ListRequest struct {
	common.PageInfo
	Type int8 `json:"type" binding:"required,oneof=1 2"`
}
type ListResponse struct {
	models.GlobalNotificationModel
	IsRead bool `json:"isRead"`
}

func (GlobalNotificationApi) ListView(c *gin.Context) {
	cr := middleware.GetBind[ListRequest](c)
	claims := jwts.GetClaims(c)
	readMsgMap := map[uint]bool{}
	query := global.DB.Where("")
	switch cr.Type {
	case 1: // 用户可见的
		// 没被用户删的
		var ugnmList []models.UserGlobalNotificationModel
		global.DB.Find(&ugnmList, "user_id = ? and is_delete = ?", claims.UserID, false)
		var msgIDList []uint
		for _, model := range ugnmList {
			readMsgMap[model.NotificationID] = model.IsRead
			msgIDList = append(msgIDList, model.ID)
		}
		query.Where("id in ?", msgIDList)
	case 2:
		if claims.Role != enum.AdminRole {
			res.FailWithMsg("权限错误", c)
			return
		}
	}
	_list, count, _ := common.ListQuery(models.GlobalNotificationModel{}, common.Options{
		PageInfo: cr.PageInfo,
		Likes:    []string{"title", "content"},
		Where:    query,
	})
	list := make([]ListResponse, 0)
	for _, model := range _list {
		list = append(list, ListResponse{
			GlobalNotificationModel: model,
			IsRead:                  readMsgMap[model.ID],
		})
	}
	res.OkWithList(list, count, c)
}
