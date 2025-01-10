package core

import (
	"fmt"
	"os"

	"github.com/LiangNing7/BlogX/conf"
	"github.com/LiangNing7/BlogX/flags"
	"gopkg.in/yaml.v2"
)

func ReadConf() (c *conf.Config) {
	byteData, err := os.ReadFile(flags.FlagOptions.File)
	if err != nil {
		panic(err)
	}
	c = new(conf.Config)
	err = yaml.Unmarshal(byteData, c)
	if err != nil {
		panic(fmt.Sprintf("yaml 配置文件格式错误 %s", err))
	}
	fmt.Printf("读取配置文件 %s 成功\n", flags.FlagOptions.File)
	return
}
