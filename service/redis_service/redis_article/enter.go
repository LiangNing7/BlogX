package redis_article

import (
	"strconv"

	"github.com/LiangNing7/BlogX/global"
)

type articleCacheType string

const (
	articleCacheLook    articleCacheType = "article_look_key"
	articleCacheDigg    articleCacheType = "article_digg_key"
	articleCacheCollect articleCacheType = "article_collect_key"
)

func set(t articleCacheType, articleID uint, increase bool) {
	num, _ := global.Redis.HGet(string(t), strconv.Itoa(int(articleID))).Int()
	if !increase {
		num--
	} else {
		num++
	}
	global.Redis.HSet(string(t), strconv.Itoa(int(articleID)), num)
}
func SetCacheLook(articleID uint, increase bool) {
	set(articleCacheLook, articleID, increase)
}
func SetCacheDigg(articleID uint, increase bool) {
	set(articleCacheDigg, articleID, increase)
}
func SetCacheCollect(articleID uint, increase bool) {
	set(articleCacheCollect, articleID, increase)
}

func get(t articleCacheType, articleID uint) int {
	num, _ := global.Redis.HGet(string(t), strconv.Itoa(int(articleID))).Int()
	return num
}
func GetCacheLook(articleID uint) int {
	return get(articleCacheLook, articleID)
}
func GetCacheDigg(articleID uint) int {
	return get(articleCacheDigg, articleID)
}
func GetCacheCollect(articleID uint) int {
	return get(articleCacheCollect, articleID)
}
func GetAll(t articleCacheType) (mps map[uint]int) {
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
func GetAllCacheLook() (mps map[uint]int) {
	return GetAll(articleCacheLook)
}
func GetAllCacheDigg() (mps map[uint]int) {
	return GetAll(articleCacheDigg)
}
func GetAllCacheCollect() (mps map[uint]int) {
	return GetAll(articleCacheCollect)
}
