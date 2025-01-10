package main

import (
	"github.com/LiangNing7/BlogX/core"
	"github.com/LiangNing7/BlogX/flags"
	"github.com/LiangNing7/BlogX/global"
	"github.com/sirupsen/logrus"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()

	logrus.Infof("123")
	logrus.Warnf("xxx")
	logrus.Debugf("yyy")
	logrus.Error("zzz")
}
