package models

import (
	"math"
	"reflect"
	"time"

	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserModel struct {
	Model
	Username       string                  `gorm:"size:32" json:"username"`
	Nickname       string                  `gorm:"size:32" json:"nickname"`
	Avatar         string                  `gorm:"size:256" json:"avatar"`
	Abstract       string                  `gorm:"size:256" json:"abstract"`
	RegisterSource enum.RegisterSourceType `json:"registerSource"` // 注册来源
	Password       string                  `gorm:"size:64" json:"-"`
	Email          string                  `gorm:"size:256" json:"email"`
	OpenID         string                  `gorm:"size:64" json:"openID"` // 第三方登陆的唯一id
	Role           enum.RoleType           `json:"role"`                  // 角色 1 管理员  2 普通用户  3 访客
	UserConfModel  *UserConfModel          `gorm:"foreignKey:UserID"  json:"-"`
	IP             string                  `json:"ip"`
	Addr           string                  `json:"addr"`
	ArticleList    []ArticleModel          `gorm:"foreignKey:UserID" json:"-"`
	LoginList      []UserLoginModel        `gorm:"foreignKey:UserID" json:"-"`
}

func (u *UserModel) GetID() uint {
	return u.ID
}

func (u *UserModel) AfterCreate(tx *gorm.DB) error {
	err := tx.Create(&UserConfModel{UserID: u.ID, OpenCollect: true, OpenFollow: true, OpenFans: true, HomeStyleID: 1}).Error
	err = tx.Create(&UserMessageConfModel{UserID: u.ID, OpenCommentMessage: true, OpenDiggMessage: true, OpenPrivateChat: true}).Error
	return err
}

func (u *UserModel) BeforeDelete(tx *gorm.DB) (err error) {
	var list = []any{
		ArticleDiggModel{},
		ArticleModel{},
		CategoryModel{},
		CollectModel{},
		CommentModel{},
		CommentDiggModel{},
		LogModel{},
		UserArticleCollectModel{},
		UserArticleLookHistoryModel{},
		UserChatActionModel{},
		UserFocusModel{},
		UserGlobalNotificationModel{},
		UserLoginModel{},
		UserTopArticleModel{},
	}
	for _, model := range list {
		count := tx.Delete(&model, "user_id = ?", u.ID).RowsAffected
		logrus.Infof("删除 %s 成功%d条", reflect.TypeOf(model).Name(), count)
	}
	var chatList []ChatModel
	tx.Find(&chatList, "send_user_id = ? or rev_user_id = ?", u.ID, u.ID).Delete(&chatList)
	logrus.Infof("删除关联对话 %d条", len(chatList))
	var messageList []MessageModel
	tx.Find(&messageList, "rev_user_id = ?", u.ID).Delete(&messageList)
	logrus.Infof("删除关联消息 %d条", len(messageList))
	return nil
}

func (u *UserModel) CodeAge() int {
	sub := time.Now().Sub(u.CreatedAt)
	return int(math.Ceil(sub.Hours() / 24 / 365))
}

type UserConfModel struct {
	UserID             uint       `gorm:"primaryKey;unique" json:"userID"`
	UserModel          UserModel  `gorm:"foreignKey:UserID" json:"-"`
	LikeTags           []string   `gorm:"type:longtext;serializer:json" json:"likeTags"`
	UpdateUsernameDate *time.Time `json:"updateUsernameDate"` // 上次修改用户名的时间
	OpenCollect        bool       `json:"openCollect"`        // 公开我的收藏
	OpenFollow         bool       `json:"openFollow"`         // 公开我的关注
	OpenFans           bool       `json:"openFans"`           // 公开我的粉丝
	HomeStyleID        uint       `json:"homeStyleID"`        // 主页样式的id
	LookCount          int        `json:"lookCount"`          // 主页的访问次数
}
