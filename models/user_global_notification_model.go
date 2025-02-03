package models

type UserGlobalNotificationModel struct {
	Model
	NotificationID uint `json:"notificationID"`
	UserID         uint `json:"userID"`
	IsRead         bool `json:"isRead"`
	IsDelete       bool `json:"isDelete"`
}
