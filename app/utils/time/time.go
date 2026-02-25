package time

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/os/gtime"
	"go.uber.org/zap"
	"strconv"
	"strings"
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

// GetStartTime start采用毫秒
func GetStartTime(start uint64, minutes int64) string {
	if start != 0 {
		t := time.UnixMilli(int64(start))
		return t.Format("2006-01-02 15:04:05")
	}

	tenMinutesAgo := time.Now().Add(-time.Minute * time.Duration(minutes))
	return tenMinutesAgo.Format("2006-01-02 15:04:05")
}

// GetEndTime end采用毫秒
func GetEndTime(end uint64, minutes int64) string {
	if end != 0 {
		t := time.UnixMilli(int64(end))
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

func SecondsToTime(seconds uint64) string {
	if seconds == 0 {
		return "00:00:00"
	}

	hours := seconds / 3600          // 总小时数
	minutes := (seconds % 3600) / 60 // 剩余分钟数
	remainingSeconds := seconds % 60 // 剩余秒数
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, remainingSeconds)
}

// TimeToString 23:59:59 => int
func TimeToString(format string) (int, error) {
	timeArr := strings.Split(format, ":")
	if len(timeArr) != 3 {
		return 0, errors.New(fmt.Sprintf("TimeToString error :%s", format))
	}

	hour, err := strconv.ParseInt(timeArr[0], 10, 64)
	if err != nil {
		zap.L().Error("TimeToString error", zap.Error(err))
		return 0, err
	}

	minute, err := strconv.ParseInt(timeArr[1], 10, 64)
	if err != nil {
		zap.L().Error("TimeToString error", zap.Error(err))
		return 0, err
	}

	second, err := strconv.ParseInt(timeArr[2], 10, 64)
	if err != nil {
		zap.L().Error("TimeToString error", zap.Error(err))
		return 0, err
	}
	ret := hour*3600 + minute*60 + second
	return int(ret), nil
}

// FormatTodayTIme 格式化当日时间 00:00 => 0
func FormatTodayTIme(time string) (int32, error) {
	arr := strings.Split(time, ":")
	if len(arr) != 2 {
		return 0, errors.New(fmt.Sprintf("日程安排错误:%s", time))
	}
	beginHour, err := strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		zap.L().Error("dealDaily error", zap.Error(err))
		return 0, err
	}

	beginMinute, err := strconv.ParseInt(arr[1], 10, 64)
	if err != nil {
		zap.L().Error("dealDaily error", zap.Error(err))
		return 0, err
	}
	beginUnix := beginHour*3600 + beginMinute*60
	return int32(beginUnix), nil
}

func SecondsToHMS(seconds int64) string {
	duration := time.Duration(seconds) * time.Second

	hours := duration / time.Hour
	duration -= hours * time.Hour

	minutes := duration / time.Minute
	duration -= minutes * time.Minute

	secondsRemaining := duration / time.Second

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secondsRemaining)
}

func FormatDateFromSecond(second int64) string {
	// 将毫秒时间戳转换为 time.Time
	t := time.UnixMilli(second)

	// 或者包含日期
	dateTimeStr := t.Format("2006-01-02 15:04:05")
	return dateTimeStr
}
