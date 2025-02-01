package article_api

import (
	"fmt"

	"github.com/LiangNing7/BlogX/common"
	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/gin-gonic/gin"
)

type CollectCreateRequest struct {
	ID       uint   `json:"id"`
	Title    string `json:"title" binding:"required,max=32"`
	Abstract string `json:"abstract"`
	Cover    string `json:"cover"`
}

func (ArticleApi) CollectCreateView(c *gin.Context) {
	cr := middleware.GetBind[CollectCreateRequest](c)
	claims := jwts.GetClaims(c)
	var model models.CollectModel
	if cr.ID == 0 {
		// 创建
		err := global.DB.Take(&model, "user_id = ? and title = ?", claims.UserID, cr.Title).Error
		if err == nil {
			res.FailWithMsg("收藏夹名称重复", c)
			return
		}
		err = global.DB.Create(&models.CollectModel{
			Title:    cr.Title,
			UserID:   claims.UserID,
			Abstract: cr.Abstract,
			Cover:    cr.Cover,
		}).Error
		if err != nil {
			res.FailWithMsg("创建收藏夹失败", c)
			return
		}
		res.OkWithMsg("创建收藏夹成功", c)
		return
	}
	err := global.DB.Take(&model, "user_id = ? and id = ?", claims.UserID, cr.ID).Error
	if err != nil {
		res.FailWithMsg("收藏夹不存在", c)
		return
	}
	err = global.DB.Model(&model).Updates(map[string]any{
		"title":    cr.Title,
		"abstract": cr.Abstract,
		"cover":    cr.Cover,
	}).Error
	if err != nil {
		res.FailWithMsg("更新收藏夹错误", c)
		return
	}
	res.OkWithMsg("更新收藏夹成功", c)
	return
}

type CollectListRequest struct {
	common.PageInfo
	UserID uint `form:"userID"`
	Type   int8 `form:"type" binding:"required,oneof=1 2 3"` // 1 查自己 2 查别人 3 后台
}
type CollectListResponse struct {
	models.CollectModel
	ArticleCount int    `json:"articleCount"`
	Nickname     string `json:"nickname,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
}

func (ArticleApi) CollectListView(c *gin.Context) {
	cr := middleware.GetBind[CollectListRequest](c)
	var preload = []string{"ArticleList"}
	switch cr.Type {
	case 1:
		claims, err := jwts.ParseTokenByGin(c)
		if err != nil {
			res.FailWithError(err, c)
			return
		}
		cr.UserID = claims.UserID
	case 2:
		var userConf models.UserConfModel
		err := global.DB.Take(&userConf, "user_id = ?", cr.UserID).Error
		if err != nil {
			res.FailWithMsg("用户不存在", c)
			return
		}
		if !userConf.OpenCollect {
			res.FailWithMsg("用户未开启我的收藏", c)
			return
		}
	case 3:
		claims, err := jwts.ParseTokenByGin(c)
		if err != nil {
			res.FailWithError(err, c)
			return
		}
		if claims.Role != enum.AdminRole {
			res.FailWithMsg("权限错误", c)
			return
		}
		preload = append(preload, "UserModel")
	}
	_list, count, _ := common.ListQuery(models.CollectModel{
		UserID: cr.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Likes:    []string{"title"},
		Preloads: preload,
	})
	var list = make([]CollectListResponse, 0)
	for _, i2 := range _list {
		list = append(list, CollectListResponse{
			CollectModel: i2,
			ArticleCount: len(i2.ArticleList),
			Nickname:     i2.UserModel.Nickname,
			Avatar:       i2.UserModel.Avatar,
		})
	}
	res.OkWithList(list, count, c)
}
func (ArticleApi) CollectRemoveView(c *gin.Context) {
	var cr = middleware.GetBind[models.RemoveRequest](c)
	var list []models.CollectModel
	query := global.DB.Where("id in ?", cr.IDList)
	claims := jwts.GetClaims(c)
	if claims.Role != enum.AdminRole {
		query.Where("user_id = ?", claims.UserID)
	}
	global.DB.Where(query).Find(&list)
	if len(list) > 0 {
		err := global.DB.Delete(&list).Error
		if err != nil {
			res.FailWithMsg("删除分类失败", c)
			return
		}
	}
	msg := fmt.Sprintf("删除收藏夹成功 共删除%d条", len(list))
	res.OkWithMsg(msg, c)
}
