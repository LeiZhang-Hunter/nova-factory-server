package order

import (
	"strconv"
	"strings"
	"time"

	"nova-factory-server/app/utils/snowflake"
)

const (
	defaultOrderPrefix = "ORD"
)

// GenerateOrderNo 使用雪花算法生成订单编号，默认前缀为 ORD。
func GenerateOrderNo() string {
	return GenerateOrderNoWithPrefix(defaultOrderPrefix)
}

// GenerateOrderNoWithPrefix 使用雪花算法生成带前缀的订单编号。
func GenerateOrderNoWithPrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		prefix = defaultOrderPrefix
	}
	now := time.Now()
	return prefix + now.Format("20060102") + strconv.FormatInt(snowflake.GenID(), 10)
}
