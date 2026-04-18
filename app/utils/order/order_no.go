package order

import (
	"fmt"
	"sync"
	"time"
)

var (
	orderNoMutex      sync.Mutex
	lastOrderNoMillis int64
	orderNoSequence   int
)

// GenerateOrderNo 生成可读且进程内有序的订单编号。
func GenerateOrderNo() string {
	now := time.Now()
	millis := now.UnixMilli()

	orderNoMutex.Lock()
	if millis == lastOrderNoMillis {
		orderNoSequence++
	} else {
		lastOrderNoMillis = millis
		orderNoSequence = 0
	}
	sequence := orderNoSequence
	orderNoMutex.Unlock()

	return fmt.Sprintf("ORD%s%03d%06d", now.Format("20060102150405"), now.Nanosecond()/int(time.Millisecond), sequence)
}
