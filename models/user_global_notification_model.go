package models

type UserGlobalNotificationModel struct {
	Model
	NotificationID uint `json:"notification_id"`
	UserID         uint `json:"user_id"`
	IsRead         bool `json:"is_read"`
	IsDelete       bool `json:"is_delete"`
}
