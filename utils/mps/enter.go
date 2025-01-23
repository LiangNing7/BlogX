package mps

import "reflect"

func StructToMap(data any, t string) (mp map[string]any) {
	mp = make(map[string]any)
	v := reflect.ValueOf(data)
	for i := 0; i < v.NumField(); i++ {
		val := v.Field(i)
		tag := v.Type().Field(i).Tag.Get(t)
		if tag == "" || tag == "-" {
			continue
		}
		if val.IsNil() {
			continue
		}
		if val.Kind() == reflect.Ptr {
			mp[tag] = val.Elem().Interface()
			continue
		}
		mp[tag] = val.Interface()
	}
	return
}
