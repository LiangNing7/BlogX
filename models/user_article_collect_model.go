package models

import "time"

type UserArticleCollectModel struct {
	UserID       uint         `gorm:"uniqueIndex:idx_name" json:"userID"`
	UserModel    UserModel    `gorm:"foreignKey:UserID" json:"-"`
	ArticleID    uint         `gorm:"uniqueIndex:idx_name" json:"articleID"`
	ArticleModel ArticleModel `gorm:"foreignKey:ArticleID" json:"-"`
	CollectID    uint         `gorm:"uniqueIndex:idx_name" json:"collectID"`    // 收藏夹的id
	CollectModel CollectModel `gorm:"foreignKey:CollectID" json:"collectModel"` // 属于哪一个收藏夹
	CreatedAt    time.Time    `json:"createdAt"`                                // 收藏的时间
}
