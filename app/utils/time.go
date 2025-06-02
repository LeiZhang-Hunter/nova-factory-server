package utils

import (
	"github.com/gogf/gf/os/gtime"
	"time"
)

func NanoToGTime(nano uint64) *gtime.Time {
	t := time.Unix(0, int64(nano/1000/1000/1000)) // 从1970-01-01 00:00:00 UTC开始计算，ns纳秒后的时间点
	return gtime.New(t)

}
