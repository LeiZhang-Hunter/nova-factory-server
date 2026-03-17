package clickhouse

import (
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type ClickHouse struct {
	db  *gorm.DB
	mtx sync.Mutex
}

func NewClickHouse() (*ClickHouse, error) {

	return &ClickHouse{db: nil}, nil
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
