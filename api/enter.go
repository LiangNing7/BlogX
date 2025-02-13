package api

import (
	"github.com/LiangNing7/BlogX/api/ai_api"
	"github.com/LiangNing7/BlogX/api/article_api"
	"github.com/LiangNing7/BlogX/api/banner_api"
	"github.com/LiangNing7/BlogX/api/captcha_api"
	"github.com/LiangNing7/BlogX/api/chat_api"
	"github.com/LiangNing7/BlogX/api/comment_api"
	"github.com/LiangNing7/BlogX/api/data_api"
	"github.com/LiangNing7/BlogX/api/focus_api"
	"github.com/LiangNing7/BlogX/api/global_notification_api"
	"github.com/LiangNing7/BlogX/api/image_api"
	"github.com/LiangNing7/BlogX/api/log_api"
	"github.com/LiangNing7/BlogX/api/search_api"
	"github.com/LiangNing7/BlogX/api/site_api"
	"github.com/LiangNing7/BlogX/api/site_msg_api"
	"github.com/LiangNing7/BlogX/api/user_api"
)

type Api struct {
	SiteApi               site_api.SiteApi
	LogApi                log_api.LogApi
	ImageApi              image_api.ImageApi
	BannerApi             banner_api.BannerApi
	CaptchaApi            captcha_api.CaptchaApi
	UserApi               user_api.UserApi
	ArticleApi            article_api.ArticleApi
	CommentApi            comment_api.CommentApi
	SiteMsgApi            site_msg_api.SiteMsgApi
	GlobalNotificationApi global_notification_api.GlobalNotificationApi
	FocusApi              focus_api.FocusApi
	ChatApi               chat_api.ChatApi
	SearchApi             search_api.SearchApi
	AiApi                 ai_api.AiApi
	DataApi               data_api.DataApi
}

var App = Api{}
