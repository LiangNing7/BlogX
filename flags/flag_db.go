package flags

import (
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/sirupsen/logrus"
)

func FlagDB() {
	err := global.DB.AutoMigrate(
		&models.UserModel{},
		&models.UserConfModel{},
		&models.ArticleModel{},
		&models.CategoryModel{},
		&models.ArticleDiggModel{},
		&models.CollectModel{},
		&models.UserArticleCollectModel{},
		&models.UserTopArticleModel{},
		&models.ImageModel{},
		&models.UserArticleLookHistoryModel{},
		&models.CommentModel{},
		&models.BannerModel{},
		&models.LogModel{},
		&models.UserLoginModel{},
		&models.GlobalNotificationModel{},
		&models.CommentDiggModel{}, // 用户点赞评论表
	)
	if err != nil {
		logrus.Errorf("数据库迁移失败 %s\n", err)
		return
	}
	logrus.Infof("数据库迁移成功!\n")
}
