package core

import (
	"github.com/LiangNing7/BlogX/global"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func InitRedis() *redis.Client {
	r := global.Config.Redis
	redisDB := redis.NewClient(&redis.Options{
		Addr:     r.Addr,     // 不写默认就是这个
		Password: r.Password, // 密码
		DB:       r.DB,       // 默认是0
	})
	_, err := redisDB.Ping().Result()
	if err != nil {
		logrus.Fatalf("redis连接失败 %s", err)
	}
	logrus.Info("redis连接成功")
	return redisDB
}
