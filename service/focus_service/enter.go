package focus_service

import (
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum/relationship_enum"
)

// CalcUserRelationship 计算好友关系
func CalcUserRelationship(A, B uint) (t relationship_enum.Relation) {
	//   2  用户2对用户1是什么关系
	var userFocusList []models.UserFocusModel
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
		return relationship_enum.RelationFans
	}
	return relationship_enum.RelationFocus
}

// CalcUserPatchRelationship 批量计算好友关系
func CalcUserPatchRelationship(A uint, BList []uint) (m map[uint]relationship_enum.Relation) {
	m = make(map[uint]relationship_enum.Relation)
	var userFocusList []models.UserFocusModel
	global.DB.Find(&userFocusList,
		"(user_id = ? and focus_user_id in ? ) or (focus_user_id = ? and user_id in ? )",
		A, BList, A, BList)
	for _, B := range BList {
		// B与A的关系
		m[B] = relationship_enum.RelationStranger
		var count int
		for _, model := range userFocusList {
			if model.FocusUserID == B {
				m[B] = relationship_enum.RelationFocus
				count++
			}
			if model.UserID == B {
				m[B] = relationship_enum.RelationFans
				count++
			}
		}
		if count == 2 {
			m[B] = relationship_enum.RelationFriend
		}
	}
	return
}
