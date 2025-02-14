package cron_service

import (
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/service/redis_service/redis_user"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SyncUser() {
	lookMap := redis_user.GetAllCacheLook()
	var list []models.UserConfModel
	global.DB.Find(&list)
	for _, model := range list {
		look := lookMap[model.UserID]
		if look == 0 {
			continue
		}
		err := global.DB.Model(&model).Updates(map[string]any{
			"look_count": gorm.Expr("look_count + ?", look),
		}).Error
		if err != nil {
			logrus.Errorf("更新失败 %s", err)
			continue
		}
		logrus.Infof("%s 更新成功", model.UserID)
	}
	// 走完之后清空掉
	redis_user.Clear()
	// 再同步回去
}
