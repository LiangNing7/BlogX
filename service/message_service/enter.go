package message_service

import (
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum/message_type_enum"
	"github.com/sirupsen/logrus"
)

// InsertCommentMessage 插入一条评论消息
func InsertCommentMessage(model models.CommentModel) {
	global.DB.Preload("UserModel").Preload("ArticleModel").Take(&model)
	err := global.DB.Create(&models.MessageModel{
		Type:               message_type_enum.CommentType,
		RevUserID:          model.ArticleModel.UserID,
		ActionUserID:       model.UserID,
		ActionUserNickname: model.UserModel.Nickname,
		ActionUserAvatar:   model.UserModel.Avatar,
		Content:            model.Content,
		ArticleID:          model.ArticleID,
		ArticleTitle:       model.ArticleModel.Title,
		CommentID:          model.ID,
	}).Error
	if err != nil {
		logrus.Error(err)
	}
}

// InsertApplyMessage 插入一条评论回复消息
func InsertApplyMessage(model models.CommentModel) {
	// TODO：回复评论的人和自己是同一个人，要不要通知？
	global.DB.Preload("ParentModel").Preload("UserModel").Preload("ArticleModel").Take(&model)
	err := global.DB.Create(&models.MessageModel{
		Type:               message_type_enum.ApplyType,
		RevUserID:          model.ParentModel.UserID,
		ActionUserID:       model.UserID,
		ActionUserNickname: model.UserModel.Nickname,
		ActionUserAvatar:   model.UserModel.Avatar,
		Content:            model.Content,
		ArticleID:          model.ArticleID,
		ArticleTitle:       model.ArticleModel.Title,
		CommentID:          model.ID,
	}).Error
	if err != nil {
		logrus.Error(err)
	}
}

// InsertDiggArticleMessage 点赞文章的消息
func InsertDiggArticleMessage(model models.ArticleDiggModel) {
	global.DB.Preload("ArticleModel").Preload("UserModel").Take(&model)
	err := global.DB.Create(&models.MessageModel{
		Type:               message_type_enum.DiggArticleType,
		RevUserID:          model.ArticleModel.UserID,
		ActionUserID:       model.UserID,
		ActionUserNickname: model.UserModel.Nickname,
		ActionUserAvatar:   model.UserModel.Avatar,
		ArticleID:          model.ArticleID,
		ArticleTitle:       model.ArticleModel.Title,
	}).Error
	if err != nil {
		logrus.Error(err)
	}
}

// InsertDiggCommentMessage 点赞评论的消息
func InsertDiggCommentMessage(model models.CommentDiggModel) {
	global.DB.Preload("CommentModel.ArticleModel").Preload("UserModel").Take(&model)
	err := global.DB.Create(&models.MessageModel{
		Type:               message_type_enum.DiggCommentType,
		RevUserID:          model.CommentModel.UserID,
		ActionUserID:       model.UserID,
		ActionUserNickname: model.UserModel.Nickname,
		ActionUserAvatar:   model.UserModel.Avatar,
		Content:            model.CommentModel.Content,
		ArticleID:          model.CommentModel.ArticleID,
		ArticleTitle:       model.CommentModel.ArticleModel.Title,
	}).Error
	if err != nil {
		logrus.Error(err)
	}
}

// InsertCollectArticleMessage 收藏文章的消息
func InsertCollectArticleMessage(model models.UserArticleCollectModel) {
	global.DB.Preload("ArticleModel").Preload("UserModel").Take(&model)
	err := global.DB.Create(&models.MessageModel{
		Type:               message_type_enum.CollectArticleType,
		RevUserID:          model.ArticleModel.UserID,
		ActionUserID:       model.UserID,
		ActionUserNickname: model.UserModel.Nickname,
		ActionUserAvatar:   model.UserModel.Avatar,
		ArticleID:          model.ArticleID,
		ArticleTitle:       model.ArticleModel.Title,
	}).Error
	if err != nil {
		logrus.Error(err)
	}
}

// InsertSystemMessage 系统消息
func InsertSystemMessage(revUserID uint, title string, content string, linkTitle string, linkHref string) {
	err := global.DB.Create(&models.MessageModel{
		Type:      message_type_enum.SystemType,
		RevUserID: revUserID,
		Title:     title,
		Content:   content,
		LinkTitle: linkTitle,
		LinkHref:  linkHref,
	}).Error
	if err != nil {
		logrus.Error(err)
	}
}
