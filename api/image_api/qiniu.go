package image_api

import (
	"fmt"

	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/service/qiniu_service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QiNiuGenTokenResponse struct {
	Token  string `json:"token"`
	Key    string `json:"key"`
	Region string `json:"region"`
	Url    string `json:"url"`
	Size   int    `json:"size"`
}

func (ImageApi) QiNiuGenToken(c *gin.Context) {
	q := global.Config.QiNiu
	if !q.Enable {
		res.FailWithMsg("未启用七牛云配置", c)
		return
	}
	token, err := qiniu_service.GenToken()
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	uid := uuid.New().String()
	key := fmt.Sprintf("%s/%s.png", q.Prefix, uid)
	url := fmt.Sprintf("%s/%s", q.Uri, key)
	res.OkWithData(QiNiuGenTokenResponse{
		Token:  token,
		Key:    key,
		Region: q.Region,
		Url:    url,
		Size:   q.Size,
	}, c)
}
