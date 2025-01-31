package image_api

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TransferDepositRequest struct {
	Url string `json:"url" binding:"required"`
}

func (ImageApi) TransferDepositView(c *gin.Context) {
	cr := middleware.GetBind[TransferDepositRequest](c)
	response, err := http.Get(cr.Url)
	if err != nil {
		res.FailWithMsg("图片请求错误", c)
		return
	}
	byteData, _ := io.ReadAll(response.Body)
	hash := utils.Md5(byteData)
	suffix := "png"
	switch response.Header.Get("Content-Type") {
	case "image/avif":
		suffix = "avif"
	}
	fmt.Println()
	filePath := fmt.Sprintf("uploads/%s/%s.%s", global.Config.Upload.UploadDir, hash, suffix)
	err = os.WriteFile(filePath, byteData, 0666)
	if err != nil {
		logrus.Error(err)
		res.FailWithMsg("图片保存失败", c)
		return
	}
	res.OkWithData("/"+filePath, c)
}
