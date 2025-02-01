package cron_service

import (
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/service/redis_service/redis_comment"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SyncComment() {
	applyMap := redis_comment.GetAllCacheApply()
	var list []models.CommentModel
	global.DB.Find(&list)
	for _, model := range list {
		apply := applyMap[model.ID]
		if apply == 0 {
			continue
		}
		err := global.DB.Model(&model).Updates(map[string]any{
			"apply_count": gorm.Expr("apply_count + ?", apply),
		}).Error
		if err != nil {
			logrus.Errorf("更新失败 %s", err)
			continue
		}
		logrus.Infof("评论%d 更新成功", model.ID)
	}
	// 走完之后清空掉
	redis_comment.Clear()
}
