package main

import (
	"fmt"

	"github.com/LiangNing7/BlogX/core"
	"github.com/LiangNing7/BlogX/flags"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/utils/jwts"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	token, err := jwts.GetToken(jwts.Claims{
		UserID: 2,
		Role:   1,
	})
	fmt.Println(token, err)
	cls, err := jwts.ParseToken("xx")
	fmt.Println(cls, err)
}
