package models

type UserLoginModel struct {
	Model
	UserID    uint      `json:"userID"`
	UserModel UserModel `gorm:"foreignKey:UserID" json:"-"` // 用户信息
	IP        string    `gorm:"size:32" json:"ip"`
	Addr      string    `gorm:"size:64" json:"addr"`
	UA        string    `gorm:"size:128" json:"ua"`
}
