package main

import (
	"fmt"

	"github.com/LiangNing7/BlogX/core"
	"github.com/LiangNing7/BlogX/flags"
)

func main() {
	flags.Parse()
	fmt.Println(flags.FlagOptions)
	core.ReadConf()
}
