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
