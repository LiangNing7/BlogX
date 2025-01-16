package log_api

import (
	"fmt"

	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/LiangNing7/BlogX/service/log_service"
	"github.com/gin-gonic/gin"
)

type LogApi struct {
}

type LogListRequest struct {
	common.PageInfo
	LogType     enum.LogType      `form:"logType"` // 日志类型 1 2 3
	Level       enum.LogLevelType `form:"level"`
	UserID      uint              `form:"userID"`
	IP          string            `form:"ip"`
	LoginStatus bool              `form:"loginStatus"`
	ServiceName string            `form:"serviceName"`
}

type LogListResponse struct {
	models.LogModel
	UserNickname string `json:"userNickname"`
	UserAvatar   string `json:"userAvatar"`
}

func (LogApi) LogListView(c *gin.Context) {
	// 分页 查询（精确查询，模糊匹配）
	var cr LogListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	list, count, err := common.ListQuery(models.LogModel{
		LogType:     cr.LogType,
		Level:       cr.Level,
		UserID:      cr.UserID,
		IP:          cr.IP,
		LoginStatus: cr.LoginStatus,
		ServiceName: cr.ServiceName,
	}, common.Options{
		PageInfo:     cr.PageInfo,
		Likes:        []string{"title"},
		Preloads:     []string{"UserModel"},
		DefaultOrder: "created_at desc",
	})

	var _list = make([]LogListResponse, 0)
	for _, logModel := range list {
		_list = append(_list, LogListResponse{
			LogModel:     logModel,
			UserNickname: logModel.UserModel.Nickname,
			UserAvatar:   logModel.UserModel.Avatar,
		})
	}

	res.OkWithList(_list, int(count), c)
	return

}

func (LogApi) LogReadView(c *gin.Context) {
	var cr models.IDRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	var log models.LogModel
	err = global.DB.Take(&log, cr.ID).Error
	if err != nil {
		res.FailWithMsg("不存在的日志", c)
		return
	}
	if !log.IsRead {
		global.DB.Model(&log).Update("is_read", true)
	}
	res.OkWithMsg("日志读取成功", c)
}

func (LogApi) LogRemoveView(c *gin.Context) {
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	log := log_service.GetLog(c)
	log.ShowRequest()
	log.ShowResponse()
	var logList []models.LogModel
	global.DB.Find(&logList, "id in ?", cr.IDList)
	if len(logList) > 0 {
		global.DB.Delete(&logList)
	}
	log.Save()
	msg := fmt.Sprintf("日志删除成功，共删除%d条", len(logList))
	res.OkWithMsg(msg, c)
}
