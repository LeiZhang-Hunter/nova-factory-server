package time

import (
	"github.com/gogf/gf/os/gtime"
	"strconv"
	"time"
)

func MicroToGTime(micro uint64) *gtime.Time {
	t := micro / 1000 / 1000
	timestamp := time.Unix(int64(t), 0)
	return gtime.New(timestamp)

}

func MillToTime(mill int64) time.Time {
	timestamp := time.UnixMilli(mill)
	return timestamp
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

// GetMonthStart 获取月份的第一天和最后一天
func GetMonthStart(myYear string, myMonth string) time.Time {
	// 数字月份必须前置补零
	if len(myMonth) == 1 {
		myMonth = "0" + myMonth
	}
	yInt, _ := strconv.Atoi(myYear)

	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, myYear+"-"+myMonth+"-01 00:00:00", loc)
	newMonth := theTime.Month()

	t1 := time.Date(yInt, newMonth, 1, 0, 0, 0, 0, time.Local)
	return t1
}

// GetWeek 获取time是周几
func GetWeek(t time.Time) int {
	return int(t.Weekday())
}
