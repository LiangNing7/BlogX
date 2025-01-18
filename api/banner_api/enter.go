package banner_api

import (
	"fmt"

	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/gin-gonic/gin"
)

type BannerApi struct {
}

type BannerCreateRequest struct {
	Cover string `json:"cover" binding:"required"`
	Href  string `json:"href"`
	Show  bool   `json:"show"`
}

func (BannerApi) BannerCreateView(c *gin.Context) {
	var cr BannerCreateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	err = global.DB.Create(&models.BannerModel{
		Cover: cr.Cover,
		Href:  cr.Href,
		Show:  cr.Show,
	}).Error
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg("添加banner成功", c)
}

type BannerListRequest struct {
	common.PageInfo
	Show bool `form:"show"`
}

func (BannerApi) BannerListView(c *gin.Context) {
	var cr BannerListRequest
	c.ShouldBindQuery(&cr)
	list, count, _ := common.ListQuery(models.BannerModel{
		Show: cr.Show,
	}, common.Options{
		PageInfo: cr.PageInfo,
	})
	res.OkWithList(list, count, c)
}
func (BannerApi) BannerRemoveView(c *gin.Context) {
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	var list []models.BannerModel
	global.DB.Find(&list, "id in ?", cr.IDList)
	if len(list) > 0 {
		global.DB.Delete(&list)
	}
	res.OkWithMsg(fmt.Sprintf("删除banner%d个，成功%d个", len(cr.IDList), len(list)), c)
}
func (BannerApi) BannerUpdateView(c *gin.Context) {
	var id models.IDRequest
	err := c.ShouldBindUri(&id)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	var cr BannerCreateRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	var model models.BannerModel
	err = global.DB.Take(&model, id.ID).Error
	if err != nil {
		res.FailWithMsg("不存在的banner", c)
		return
	}
	err = global.DB.Model(&model).Updates(map[string]any{
		"cover": cr.Cover,
		"href":  cr.Href,
		"show":  cr.Show,
	}).Error
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg("banner更新成功", c)
}
