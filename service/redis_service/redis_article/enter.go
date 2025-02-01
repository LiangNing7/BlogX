package redis_article

import (
	"fmt"
	"strconv"

	"github.com/LiangNing7/BlogX/global"
	"github.com/LiangNing7/BlogX/utils/date"
	"github.com/sirupsen/logrus"
)

type articleCacheType string

const (
	articleCacheLook    articleCacheType = "article_look_key"
	articleCacheDigg    articleCacheType = "article_digg_key"
	articleCacheCollect articleCacheType = "article_collect_key"
	articleCacheComment articleCacheType = "article_comment_key"
)

func set(t articleCacheType, articleID uint, n int) {
	num, _ := global.Redis.HGet(string(t), strconv.Itoa(int(articleID))).Int()
	num += n
	global.Redis.HSet(string(t), strconv.Itoa(int(articleID)), num)
}

func SetCacheLook(articleID uint, increase bool) {
	var n = 1
	if !increase {
		n = -1
	}
	set(articleCacheLook, articleID, n)
}
func SetCacheDigg(articleID uint, increase bool) {
	var n = 1
	if !increase {
		n = -1
	}
	set(articleCacheDigg, articleID, n)
}
func SetCacheCollect(articleID uint, increase bool) {
	var n = 1
	if !increase {
		n = -1
	}
	set(articleCacheCollect, articleID, n)
}
func SetCacheComment(articleID uint, n int) {
	set(articleCacheComment, articleID, n)
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
func GetCacheComment(articleID uint) int {
	return get(articleCacheComment, articleID)
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
func GetAllCacheComment() (mps map[uint]int) {
	return GetAll(articleCacheComment)
}

func Clear() {
	err := global.Redis.Del("article_look_key", "article_digg_key", "article_collect_key").Err()
	if err != nil {
		logrus.Error(err)
	}
}

func SetUserArticleHistoryCache(articleID, userID uint) {
	key := fmt.Sprintf("histroy_%d", userID)
	field := fmt.Sprintf("%d", articleID)

	endTime := date.GetNowAfter()
	err := global.Redis.HSet(key, field, "").Err()
	if err != nil {
		logrus.Error(err)
		return
	}
	err = global.Redis.ExpireAt(key, endTime).Err()
	if err != nil {
		logrus.Error(err)
		return
	}
}
func GetUserArticleHistoryCache(articleID, userID uint) (ok bool) {
	key := fmt.Sprintf("histroy_%d", userID)
	field := fmt.Sprintf("%d", articleID)
	err := global.Redis.HGet(key, field).Err()
	if err != nil {
		return false
	}
	return true
}
