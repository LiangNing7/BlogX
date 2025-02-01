package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func InList[T comparable](key T, list []T) bool {
	for _, s := range list {
		if key == s {
			return true
		}
	}
	return false
}

func Md5(data []byte) string {
	md5New := md5.New()
	md5New.Write(data)
	return hex.EncodeToString(md5New.Sum(nil))
}

// Unique 切片去重升级版 泛型参数 利用map的key不能重复的特性+append函数  一次for循环搞定
func Unique[T comparable](ss []T) []T {
	size := len(ss)
	if size == 0 {
		return []T{}
	}
	newSlices := make([]T, 0) // 这里新建一个切片,大于为0, 因为我们不知道有几个非重复数据,后面都使用append来动态增加并扩容
	m1 := make(map[T]bool)
	for _, v := range ss {
		if _, ok := m1[v]; !ok { // 如果数据不在map中,放入
			m1[v] = true                     // 保存到map中,用于下次判断
			newSlices = append(newSlices, v) // 将数据放入新的切片中
		}
	}
	return newSlices
}
