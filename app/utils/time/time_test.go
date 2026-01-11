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
