package redis_article

import (
	"fmt"
	"strconv"
	"time"

	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/utils/date"
	"github.com/sirupsen/logrus"
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

func SetUserArticleHistoryCache(articleID, userID uint) {
	key := fmt.Sprintf("histroy_%d_%d", userID, articleID)
	now := time.Now().Local()
	endTime := date.GetNowAfter()
	subTime := endTime.Sub(now)
	err := global.Redis.Set(key, "", subTime).Err()
	if err != nil {
		logrus.Error(err)
		return
	}
}
func GetUserArticleHistoryCache(articleID, userID uint) (ok bool) {
	key := fmt.Sprintf("histroy_%d_%d", userID, articleID)
	err := global.Redis.Get(key).Err()
	if err != nil {
		return false
	}
	return true
}
