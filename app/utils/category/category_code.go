package category

import (
	"strconv"
	"time"

	"nova-factory-server/app/utils/snowflake"
)

const (
	defaultProductCategoryCodePrefix = "PCAT"
)

// GenerateProductCategoryCode 生成 ERP 产品分类编码。
func GenerateProductCategoryCode() string {
	now := time.Now()
	return defaultProductCategoryCodePrefix + now.Format("20060102150405") + strconv.FormatInt(snowflake.GenID(), 10)
}
