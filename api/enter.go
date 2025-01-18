package api

import (
	"github.com/LiangNing7/BlogX/api/banner_api"
	"github.com/LiangNing7/BlogX/api/image_api"
	"github.com/LiangNing7/BlogX/api/log_api"
	"github.com/LiangNing7/BlogX/api/site_api"
)

type Api struct {
	SiteApi   site_api.SiteApi
	LogApi    log_api.LogApi
	ImageApi  image_api.ImageApi
	BannerApi banner_api.BannerApi
}

var App = Api{}
