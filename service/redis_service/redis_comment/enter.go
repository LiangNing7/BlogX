package redis_comment

import (
	"strconv"

	"github.com/LiangNing7/BlogX/global"
	"github.com/sirupsen/logrus"
)

type commentCacheType string

const (
	commentCacheApply commentCacheType = "comment_apply_key"
)

func set(t commentCacheType, commentID uint, n int) {
	num, _ := global.Redis.HGet(string(t), strconv.Itoa(int(commentID))).Int()
	num += n
	global.Redis.HSet(string(t), strconv.Itoa(int(commentID)), num)
}
func SetCacheApply(commentID uint, n int) {
	set(commentCacheApply, commentID, n)
}
func get(t commentCacheType, commentID uint) int {
	num, _ := global.Redis.HGet(string(t), strconv.Itoa(int(commentID))).Int()
	return num
}
func GetCacheApply(commentID uint) int {
	return get(commentCacheApply, commentID)
}
func GetAll(t commentCacheType) (mps map[uint]int) {
	res, err := global.Redis.HGetAll(string(t)).Result()
	if err != nil {
		return
	}
	mps = make(map[uint]int)
	for key, numS := range res {
		iK, err := strconv.Atoi(key)
		if err != nil {
			continue
		}
		iN, err := strconv.Atoi(numS)
		if err != nil {
			continue
		}
		mps[uint(iK)] = iN
	}
	return mps
}
func GetAllCacheApply() (mps map[uint]int) {
	return GetAll(commentCacheApply)
}

func Clear() {
	err := global.Redis.Del("comment_apply_key").Err()
	if err != nil {
		logrus.Error(err)
	}
}
