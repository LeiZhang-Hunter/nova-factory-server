package clickhouse

import (
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

var clickhouseOnce sync.Once
var clickhouseInstance *ClickHouse

type ClickHouse struct {
	db  *gorm.DB
	mtx sync.Mutex
}

func newClickHouse() *ClickHouse {

	return &ClickHouse{db: nil}
}

func GetClickHouse() *ClickHouse {
	if clickhouseInstance != nil {
		return clickhouseInstance
	}
	clickhouseOnce.Do(func() {
		clickhouseInstance = newClickHouse()
	})
	return clickhouseInstance
}

func (c *ClickHouse) DB() *gorm.DB {
	if c.db == nil {
		c.mtx.Lock()
		defer c.mtx.Unlock()
		if c.db != nil {
			return c.db
		}
		dsn := viper.GetString("clickhouse.link")
		db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
		if err != nil {
			zap.L().Error("connect clickhouse fail", zap.Error(err))
			return c.db
		}
		c.db = db
	}
	return c.db
}
