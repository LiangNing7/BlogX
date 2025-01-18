package main

import (
	"context"
	"fmt"

	"github.com/LiangNing7/BlogX/core"
	"github.com/LiangNing7/BlogX/flags"
	"github.com/LiangNing7/BlogX/global"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
)

func SendFile(file string) (err error) {
	mac := credentials.NewCredentials(global.Config.QiNiu.AccessKey, global.Config.QiNiu.SecretKey)
	key := "blogx/jk雷神.jpg"
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	err = uploadManager.UploadFile(context.Background(), file, &uploader.ObjectOptions{
		BucketName: global.Config.QiNiu.Bucket,
		ObjectName: &key,
		FileName:   "xx.jpg",
	}, nil)
	return err
}

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	err := SendFile("uploads/images/jk雷神.jpg")
	fmt.Println(err)
}
