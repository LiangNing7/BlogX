package common

import (
	"reflect"

	"github.com/LiangNing7/BlogX/global"
	"gorm.io/gorm"
)

type ModelMap interface {
	GetID() uint
}
type ScanOption struct {
	Where *gorm.DB
}

func ScanMap[T ModelMap](model T, option ScanOption) (mp map[uint]T) {
	var list []T
	query := global.DB.Where(model)
	if option.Where != nil {
		query = query.Where(option.Where)
	}
	query.Find(&list)
	mp = map[uint]T{}
	for _, m := range list {
		mp[m.GetID()] = m
	}
	return
}
func ScanMapV2[T any](model T, option ScanOption) (mp map[uint]T) {
	var list []T
	query := global.DB.Where(model)
	if option.Where != nil {
		query = query.Where(option.Where)
	}
	query.Find(&list)
	mp = map[uint]T{}
	for _, m := range list {
		v := reflect.ValueOf(m)
		idField := v.FieldByName("ID")
		uid, ok := idField.Interface().(uint)
		if !ok {
			continue
		}
		mp[uid] = m
	}
	return
}
