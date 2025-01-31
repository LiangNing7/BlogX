package api

import (
	"github.com/LiangNing7/BlogX/api/article_api"
	"github.com/LiangNing7/BlogX/api/banner_api"
	"github.com/LiangNing7/BlogX/api/captcha_api"
	"github.com/LiangNing7/BlogX/api/image_api"
	"github.com/LiangNing7/BlogX/api/log_api"
	"github.com/LiangNing7/BlogX/api/site_api"
	"github.com/LiangNing7/BlogX/api/user_api"
)

type Api struct {
	SiteApi    site_api.SiteApi
	LogApi     log_api.LogApi
	ImageApi   image_api.ImageApi
	BannerApi  banner_api.BannerApi
	CaptchaApi captcha_api.CaptchaApi
	UserApi    user_api.UserApi
	ArticleApi article_api.ArticleApi
}

var App = Api{}
