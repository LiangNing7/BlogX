package global

import (
	"github.com/LiangNing7/BlogX/conf"
	"gorm.io/gorm"
)

var (
	Config *conf.Config
	DB     *gorm.DB
)
