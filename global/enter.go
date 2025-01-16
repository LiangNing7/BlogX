package global

import (
	"github.com/LiangNing7/BlogX/conf"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	Config *conf.Config
	DB     *gorm.DB
	Redis  *redis.Client
)
