package core

import (
	"github.com/LiangNing7/BlogX/global"
	river "github.com/LiangNing7/BlogX/service/river_service"
	"github.com/sirupsen/logrus"
)

func InitMysqlES() {
	if !global.Config.River.Enable {
		logrus.Infof("关闭mysql同步操作")
		return
	}
	r, err := river.NewRiver()
	if err != nil {
		logrus.Fatal(err)
	}
	go r.Run()
}
