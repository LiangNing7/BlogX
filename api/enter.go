package api

import (
	"github.com/LiangNing7/BlogX/api/log_api"
	"github.com/LiangNing7/BlogX/api/site_api"
)

type Api struct {
	SiteApi site_api.SiteApi
	LogApi  log_api.LogApi
}

var App = Api{}
