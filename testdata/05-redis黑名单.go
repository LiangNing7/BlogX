package main

import (
	"fmt"

	"github.com/LiangNing7/BlogX/core"
	"github.com/LiangNing7/BlogX/flags"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/service/redis_service/redis_jwt"
	"github.com/LiangNing7/BlogX/utils/jwts"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.Redis = core.InitRedis()
	token, err := jwts.GetToken(jwts.Claims{
		UserID: 2,
		Role:   1,
	})
	fmt.Println(token, err)
	redis_jwt.TokenBlack(token, redis_jwt.UserBlackType)
	blk, ok := redis_jwt.HasTokenBlack(token)
	fmt.Println(blk, ok)
}
