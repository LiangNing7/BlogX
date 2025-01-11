package api

import "github.com/LiangNing7/BlogX/api/site_api"

type Api struct {
	SiteApi site_api.SiteApi
}

var App = Api{}
