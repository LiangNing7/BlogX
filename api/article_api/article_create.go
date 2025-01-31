package article_api

import (
	"bytes"

	"github.com/LiangNing7/BlogX/common/res"
	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/middleware"
	"github.com/LiangNing7/BlogX/models"
	"github.com/LiangNing7/BlogX/models/ctype"
	"github.com/LiangNing7/BlogX/models/enum"
	"github.com/LiangNing7/BlogX/utils/jwts"
	"github.com/LiangNing7/BlogX/utils/markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

type ArticleCreateRequest struct {
	Title       string             `json:"title" binding:"required"`
	Abstract    string             `json:"abstract"`
	Content     string             `json:"content" binding:"required"`
	CategoryID  *uint              `json:"categoryID"`
	TagList     ctype.List         `json:"tagList"`
	Cover       string             `json:"cover"`
	OpenComment bool               `json:"openComment"`
	Status      enum.ArticleStatus `json:"status" binding:"required,oneof=1 2"`
}

func (ArticleApi) ArticleCreateView(c *gin.Context) {
	cr := middleware.GetBind[ArticleCreateRequest](c)
	user, err := jwts.GetClaims(c).GetUser()
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}
	// 判断分类id是不是自己创建的
	var category models.CategoryModel
	if cr.CategoryID != nil {
		err = global.DB.Take(&category, "id = ? and user_id = ?", *cr.CategoryID, user.ID).Error
		if err != nil {
			res.FailWithMsg("文章分类不存在", c)
			return
		}
	}
	// 文章正文防xss注入
	contentDoc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(cr.Content)))
	if err != nil {
		res.FailWithMsg("正文解析错误", c)
		return
	}
	contentDoc.Find("script").Remove()
	contentDoc.Find("img").Remove()
	contentDoc.Find("iframe").Remove()
	cr.Content = contentDoc.Text()
	// 如果不传简介，那么从正文中取前30个字符
	if cr.Abstract == "" {
		// 把markdown转成html，再取文本
		html := markdown.MdToHtml(cr.Content)
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(html)))
		if err != nil {
			res.FailWithMsg("正文解析错误", c)
			return
		}
		htmlText := doc.Text()
		cr.Abstract = htmlText
		if len(htmlText) > 200 {
			// 如果大于200，就取前200
			cr.Abstract = string([]rune(htmlText)[:200])
		}
	}
	// 正文内容图片转存
	// 1.图片过多，同步做，接口耗时高  异步做，

	var article = models.ArticleModel{
		Title:       cr.Title,
		Abstract:    cr.Abstract,
		Content:     cr.Content,
		UserID:      user.ID,
		TagList:     cr.TagList,
		Cover:       cr.Cover,
		OpenComment: cr.OpenComment,
		CategoryID:  cr.CategoryID,
		Status:      cr.Status,
	}
	if cr.Status == enum.ArticleStatusExamine && global.Config.Site.Article.NoExamine {
		article.Status = enum.ArticleStatusPublished
	}
	err = global.DB.Create(&article).Error
	if err != nil {
		res.FailWithMsg("文章创建失败", c)
		return
	}
	res.OkWithMsg("文章创建成功", c)
}
