package date

import "time"

func GetNowAfter() time.Time {
	// 获取当前时间
	now := time.Now()
	// 获取当前时区
	location := time.Local
	// 设置今天的结束时间为23:59:59，基于当前时区
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, location)
	return endTime
}
