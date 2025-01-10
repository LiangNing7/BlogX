package models

type CollectModel struct {
	Model
	Title        string    `gorm:"size:32" json:"title"`
	Abstract     string    `gorm:"size:256" json:"abstract"`
	Cover        string    `gorm:"size:256" json:"cover"`
	ArticleCount int       `json:"articleCount"`
	UserID       uint      `json:"userID"`
	UserModel    UserModel `gorm:"foreignKey:UserID" json:"-"`
}
