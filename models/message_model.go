package models

import "github.com/LiangNing7/BlogX/models/enum/message_type_enum"

type MessageModel struct {
	Model
	Type               message_type_enum.Type `json:"type"`
	RevUserID          uint                   `json:"revUserID"` // 接收人的id
	ActionUserID       uint                   `json:"ActionUserID"`
	ActionUserNickname string                 `json:"actionUserNickname"`
	ActionUserAvatar   string                 `json:"actionUserAvatar"`
	Title              string                 `json:"title"`
	Content            string                 `json:"content"`
	ArticleID          uint                   `json:"articleID"`
	ArticleTitle       string                 `json:"articleTitle"`
	CommentID          uint                   `json:"commentID"`
	LinkTitle          string                 `json:"linkTitle"`
	LinkHref           string                 `json:"linkHref"`
	IsRead             bool                   `json:"isRead"`
}
