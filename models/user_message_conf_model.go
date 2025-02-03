package models

type UserMessageConfModel struct {
	UserID             uint      `gorm:"primaryKey;unique" json:"userID"`
	UserModel          UserModel `gorm:"foreignKey:UserID" json:"-"`
	OpenCommentMessage bool      `json:"openCommentMessage"` // 开启回复和评论
	OpenDiggMessage    bool      `json:"openDiggMessage"`    // 开启赞和收藏
	OpenPrivateChat    bool      `json:"openPrivateChat"`    // 是否开启私聊
}
