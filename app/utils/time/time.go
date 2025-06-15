package time

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

func GetStartTime(start uint64, minutes int64) string {
	if start != 0 {
		t := time.Unix(0, int64(start)*1e6)
		return t.Format("2006-01-02 15:04:05")
	}

	tenMinutesAgo := time.Now().Add(-time.Minute * time.Duration(minutes))
	return tenMinutesAgo.Format("2006-01-02 15:04:05")
}

func GetEndTime(end uint64, minutes int64) string {
	if end != 0 {
		t := time.Unix(0, int64(end)*1e6)
		return t.Format("2006-01-02 15:04:05")
	}
	tenMinutesAgo := time.Unix(0, int64(end)*1e6).Add(-time.Minute * time.Duration(minutes))
	return tenMinutesAgo.Format("2006-01-02 15:04:05")
}

func GetEndTimeUseNow(end uint64, useNow bool) string {
	if end != 0 {
		t := time.Unix(0, int64(end)*1e6)
		return t.Format("2006-01-02 15:04:05")
	} else {
		t := time.Now()
		return t.Format("2006-01-02 15:04:05")
	}
}
