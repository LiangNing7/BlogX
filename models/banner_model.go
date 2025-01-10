package models

type BannerModel struct {
	Model
	Cover string `json:"cover"` // 图片链接
	Href  string `json:"href"`  // 跳转链接
}
