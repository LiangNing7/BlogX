package core

import (
	"fmt"
	"os"

	"github.com/LiangNing7/BlogX/flags"
	"gopkg.in/yaml.v2"
)

type System struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type Config struct {
	System System `yaml:"system"`
}

func ReadConf() {
	byteData, err := os.ReadFile(flags.FlagOptions.File)
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(byteData, &config)
	if err != nil {
		panic(fmt.Sprintf("yaml 配置文件格式错误 %s", err))
	}
	fmt.Printf("读取配置文件 %s 成功", flags.FlagOptions.File)
}
