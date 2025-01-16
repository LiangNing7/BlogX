package site_api

import (
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/service/log_service"
	"github.com/gin-gonic/gin"
)

type SiteApi struct {
}

func (SiteApi) SiteInfoView(c *gin.Context) {
	res.OkWithData("xx", c)
	return
}

type SiteUpdateRequest struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required" label:"年龄"`
}

func (SiteApi) SiteUpdateView(c *gin.Context) {
	var cr SiteUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	log := log_service.NewActionLogByGin(c)
	log.SetItem("name", cr.Name)
	log.ShowResponseHeader()
	log.ShowResponse()

	res.OkWithMsg("更新成功", c)
	return
}
