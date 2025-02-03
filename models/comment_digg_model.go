package models

import "time"

type CommentDiggModel struct {
	Model
	UserID       uint         `gorm:"uniqueIndex:idx_name" json:"userID"`
	UserModel    UserModel    `gorm:"foreignKey:UserID" json:"-"`
	CommentID    uint         `gorm:"uniqueIndex:idx_name" json:"commentID"`
	CommentModel CommentModel `gorm:"foreignKey:CommentID" json:"-"`
	CreatedAt    time.Time    `json:"createdAt"` // 点赞的时间
}
