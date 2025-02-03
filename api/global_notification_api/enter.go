package global_notification_api

import (
	"fmt"

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
	Type int8 `form:"type" binding:"required,oneof=1 2"`
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
		global.DB.Find(&ugnmList, "user_id = ?", claims.UserID)

		var msgIDList []uint
		for _, model := range ugnmList {
			if model.IsDelete {
				msgIDList = append(msgIDList, model.ID)
				continue
			}
			if model.IsRead {
				readMsgMap[model.NotificationID] = true
			}
		}
		if len(msgIDList) > 0 {
			query.Where("id not in ?", msgIDList)
		}

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

func (GlobalNotificationApi) RemoveAdminView(c *gin.Context) {
	cr := middleware.GetBind[models.RemoveRequest](c)

	var list []models.GlobalNotificationModel
	global.DB.Find(&list, "id in ?", cr.IDList)

	if len(list) > 0 {
		global.DB.Delete(&list)
	}

	res.OkWithMsg(fmt.Sprintf("删除%d条全局消息，成功%d个", len(cr.IDList), len(list)), c)
}

type UserMsgActionRequest struct {
	ID   uint `json:"id" binding:"required"`
	Type int8 `json:"type" binding:"required,oneof=1 2"` // 1 读取 2 删除
}

// UserMsgActionView 用户读取或者用户删除全局消息
func (GlobalNotificationApi) UserMsgActionView(c *gin.Context) {
	cr := middleware.GetBind[UserMsgActionRequest](c)

	var msg models.GlobalNotificationModel
	err := global.DB.Take(&msg, cr.ID).Error
	if err != nil {
		res.FailWithMsg("消息不存在", c)
		return
	}

	claims := jwts.GetClaims(c)

	model := models.UserGlobalNotificationModel{
		NotificationID: cr.ID,
		UserID:         claims.UserID,
	}
	m := "消息读取成功"
	if cr.Type == 1 {
		model.IsRead = true
	} else {
		model.IsDelete = true
		m = "消息删除成功"
	}

	// 看一看之前有没有操作过
	var ugnm models.UserGlobalNotificationModel
	err = global.DB.Take(&ugnm, "notification_id = ? and user_id = ?", cr.ID, claims.UserID).Error
	// 之前这个用户对这个消息没有操作过
	// 之前对这个消息有读取操作
	// 之前对这个消息有删除操作
	// 先删除再读取
	if err != nil {
		global.DB.Create(&model)
		res.OkWithMsg("消息读取成功", c)
		return
	}
	if ugnm.IsDelete {
		res.FailWithMsg("消息已删除", c)
		return
	}
	if ugnm.IsRead {
		// 如果现在是删除操作，那就更新
		if model.IsDelete {
			global.DB.Model(&ugnm).Update("is_delete", true)
		}
	}

	res.OkWithMsg(m, c)

}
