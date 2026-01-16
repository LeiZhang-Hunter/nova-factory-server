package time

import (
	"fmt"
	"testing"
	systime "time"
)

func TestFormatDateFromSecond(t *testing.T) {
	end := systime.Now().UnixMilli()
	start := end - 7*86400*1000
	startTime := FormatDateFromSecond(start)
	endTime := FormatDateFromSecond(end)
	fmt.Println(startTime, endTime)
}

func TestStartTime(t *testing.T) {
	start := systime.Now().UnixMilli()
	startStr := GetStartTime(uint64(start), 0)
	fmt.Println(startStr)

	start = systime.Now().UnixMilli() - 86400*1000
	startStr = GetStartTime(uint64(start), 0)
	fmt.Println(startStr)
}

func TestEndTime(t *testing.T) {
	start := systime.Now().UnixMilli()
	startStr := GetEndTime(uint64(start), 0)
	fmt.Println(startStr)

	start = systime.Now().UnixMilli() - 86400*1000
	startStr = GetEndTime(uint64(start), 0)
	fmt.Println(startStr)
}
