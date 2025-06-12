package utils

import (
	"fmt"
	"github.com/gogf/gf/os/gtime"
	"time"
)

func MicroToGTime(micro uint64) *gtime.Time {
	fmt.Println(time.Now().Unix())
	t := micro / 1000 / 1000
	timestamp := time.Unix(int64(t), 0)
	return gtime.New(timestamp)

}
