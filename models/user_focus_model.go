package models

import (
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models/enum/relationship_enum"
)

type UserFocusModel struct {
	Model
	UserID         uint      `json:"userID"` // 用户id
	UserModel      UserModel `gorm:"foreignKey:UserID" json:"-"`
	FocusUserID    uint      `json:"focusUserID"` // 关注的用户
	FocusUserModel UserModel `gorm:"foreignKey:FocusUserID" json:"-"`
}

// CalcUserRelationship 计算好友关系
func CalcUserRelationship(A, B uint) (t relationship_enum.Relation) {
	//   2  用户2对用户1是什么关系
	var userFocusList []UserFocusModel
	global.DB.Find(&userFocusList,
		"(user_id = ? and focus_user_id = ? ) or (focus_user_id = ? and user_id = ? )",
		A, B, A, B)
	if len(userFocusList) == 2 {
		return relationship_enum.RelationFriend
	}
	if len(userFocusList) == 0 {
		return relationship_enum.RelationStranger
	}
	focus := userFocusList[0]
	if focus.FocusUserID == A {
		return relationship_enum.RelationFocus
	}
	return relationship_enum.RelationFans
}
