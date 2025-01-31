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

type CategoryCreateRequest struct {
	ID    uint   `json:"id"`
	Title string `json:"title" binding:"required,max=32"`
}

func (ArticleApi) CategoryCreateView(c *gin.Context) {
	cr := middleware.GetBind[CategoryCreateRequest](c)
	claims := jwts.GetClaims(c)
	var model models.CategoryModel
	if cr.ID == 0 {
		// 创建
		err := global.DB.Take(&model, "user_id = ? and title = ?", claims.UserID, cr.Title).Error
		if err == nil {
			res.FailWithMsg("分类名称重复", c)
			return
		}
		err = global.DB.Create(&models.CategoryModel{
			Title:  cr.Title,
			UserID: claims.UserID,
		}).Error
		if err != nil {
			res.FailWithMsg("创建分类错误", c)
			return
		}
		res.OkWithMsg("创建分类成功", c)
		return
	}
	err := global.DB.Take(&model, "user_id = ? and id = ?", claims.UserID, cr.ID).Error
	if err != nil {
		res.FailWithMsg("分类不存在", c)
		return
	}
	err = global.DB.Model(&model).Update("title", cr.Title).Error
	if err != nil {
		res.FailWithMsg("更新分类错误", c)
		return
	}
	res.OkWithMsg("更新分类成功", c)
	return
}

type CategoryListRequest struct {
	common.PageInfo
	UserID uint `form:"userID"`
	Type   int8 `form:"type" binding:"required,oneof=1 2 3"` // 1 查自己 2 查别人 3 后台
}
type CategoryListResponse struct {
	models.CategoryModel
	ArticleCount int    `json:"articleCount"`
	Nickname     string `json:"nickname,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
}

func (ArticleApi) CategoryListView(c *gin.Context) {
	cr := middleware.GetBind[CategoryListRequest](c)
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
	_list, count, _ := common.ListQuery(models.CategoryModel{
		UserID: cr.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Likes:    []string{"title"},
		Preloads: preload,
	})
	var list = make([]CategoryListResponse, 0)
	for _, i2 := range _list {
		list = append(list, CategoryListResponse{
			CategoryModel: i2,
			ArticleCount:  len(i2.ArticleList),
			Nickname:      i2.UserModel.Nickname,
			Avatar:        i2.UserModel.Avatar,
		})
	}
	res.OkWithList(list, count, c)
}

func (ArticleApi) CategoryRemoveView(c *gin.Context) {
	var cr = middleware.GetBind[models.RemoveRequest](c)
	var list []models.CategoryModel
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
	msg := fmt.Sprintf("删除分类成功 共删除%d条", len(list))
	res.OkWithMsg(msg, c)
}
