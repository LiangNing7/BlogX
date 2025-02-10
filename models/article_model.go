package models

import (
	_ "embed"

	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models/ctype"
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/LiangNing7/BlogX/service/text_service"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ArticleModel struct {
	Model
	Title         string             `gorm:"size:32" json:"title"`
	Abstract      string             `gorm:"size:256" json:"abstract"`
	Content       string             `json:"content,omitempty"`
	CategoryID    *uint              `json:"categoryID"` // 分类的id
	CategoryModel *CategoryModel     `gorm:"foreignKey:CategoryID" json:"-"`
	TagList       ctype.List         `gorm:"type:longtext" json:"tagList"` // 标签列表
	Cover         string             `gorm:"size:256" json:"cover"`
	UserID        uint               `json:"userID"`
	UserModel     UserModel          `gorm:"foreignKey:UserID" json:"-"`
	LookCount     int                `json:"lookCount"`
	DiggCount     int                `json:"diggCount"`
	CommentCount  int                `json:"commentCount"`
	CollectCount  int                `json:"collectCount"`
	OpenComment   bool               `json:"openComment"` // 开启评论
	Status        enum.ArticleStatus `json:"status"`      // 状态 草稿 审核中  已发布
}

//go:embed mappings/article_mapping.json
var articleMapping string

func (ArticleModel) Mapping() string {
	return articleMapping
}
func (ArticleModel) Index() string {
	return "article_index"
}

func (a *ArticleModel) BeforeDelete(tx *gorm.DB) (err error) {
	// 评论
	var commentList []CommentModel
	global.DB.Find(&commentList, "article_id = ?", a.ID).Delete(&commentList)
	// 点赞
	var diggList []ArticleDiggModel
	global.DB.Find(&diggList, "article_id = ?", a.ID).Delete(&diggList)
	// 收藏
	var collectList []UserArticleCollectModel
	global.DB.Find(&collectList, "article_id = ?", a.ID).Delete(&collectList)
	// 置顶
	var topList []UserTopArticleModel
	global.DB.Find(&topList, "article_id = ?", a.ID).Delete(&topList)
	// 浏览
	var lookList []UserArticleLookHistoryModel
	global.DB.Find(&lookList, "article_id = ?", a.ID).Delete(&lookList)
	logrus.Infof("删除关联评论 %d 条", len(commentList))
	logrus.Infof("删除关联点赞 %d 条", len(diggList))
	logrus.Infof("删除关联收藏 %d 条", len(collectList))
	logrus.Infof("删除关联置顶 %d 条", len(topList))
	logrus.Infof("删除关联浏览 %d 条", len(lookList))
	return
}

func (a *ArticleModel) AfterCreate(tx *gorm.DB) (err error) {
	// 创建文章之后的钩子函数
	// 只有发布中的文章会放到全文搜索里面去
	if a.Status != enum.ArticleStatusPublished {
		return nil
	}
	textList := text_service.MdContentTransformation(a.ID, a.Title, a.Content)
	var list []TextModel
	for _, model := range textList {
		list = append(list, TextModel{
			ArticleID: model.ArticleID,
			Head:      model.Head,
			Body:      model.Body,
		})
	}
	err = tx.Create(&list).Error
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return nil
}
func (a *ArticleModel) AfterDelete(tx *gorm.DB) (err error) {
	// 删除之后
	var textList []TextModel
	tx.Find(&textList, "article_id = ?", a.ID)
	if len(textList) > 0 {
		logrus.Infof("删除全文记录 %d", len(textList))
		tx.Delete(&textList)
	}
	return nil
}

func (a *ArticleModel) AfterUpdate(tx *gorm.DB) (err error) {
	// 正文发生了变化，才去做转换
	a.AfterDelete(tx)
	a.AfterCreate(tx)
	return nil
}
