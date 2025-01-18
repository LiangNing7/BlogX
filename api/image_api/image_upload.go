package image_api

import (
	"fmt"
	"io"

	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/utils"
	file2 "github.com/LiangNing7/BlogX/utils/file"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (ImageApi) ImageUploadView(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	// 文件大小判断
	s := global.Config.Upload.Size
	if fileHeader.Size > s*1024*1024 {
		res.FailWithMsg(fmt.Sprintf("文件大小大于%dMB", s), c)
		return
	}
	// 后缀判断
	filename := fileHeader.Filename
	suffix, err := file2.ImageSuffixJudge(filename)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	// 文件hash
	file, err := fileHeader.Open()
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	byteData, _ := io.ReadAll(file)
	hash := utils.Md5(byteData)
	// 判断这个hash有没有
	var model models.ImageModel
	err = global.DB.Take(&model, "hash = ?", hash).Error
	if err == nil {
		// 找到了
		logrus.Infof("上传图片重复 %s <==> %s  %s", filename, model.Filename, hash)
		res.Ok(model.WebPath(), "上传成功", c)
		return
	}
	// 文件名称一样，但是文件内容不一样
	filePath := fmt.Sprintf("uploads/%s/%s.%s", global.Config.Upload.UploadDir, hash, suffix)
	// 入库
	model = models.ImageModel{
		Filename: filename,
		Path:     filePath,
		Size:     fileHeader.Size,
		Hash:     hash,
	}
	err = global.DB.Create(&model).Error
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	c.SaveUploadedFile(fileHeader, filePath)
	res.Ok(model.WebPath(), "图片上传成功", c)
}
