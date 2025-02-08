package models

type UserChatActionModel struct {
	Model
	UserID   uint `json:"userID"`
	ChatID   uint `json:"chatID"`
	IsRead   bool `json:"isRead"`
	IsDelete bool `json:"isDelete"`
}
