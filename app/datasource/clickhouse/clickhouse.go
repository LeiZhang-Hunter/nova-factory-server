package clickhouse

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"sync"
)

type ClickHouse struct {
	db  *gorm.DB
	mtx sync.Mutex
}

func NewClickHouse() (*ClickHouse, error) {
	dsn := viper.GetString("clickhouse.link")
	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("connect clickhouse fail", zap.Error(err))
		return &ClickHouse{
			db: nil,
		}, nil
	}
	return &ClickHouse{db: db}, nil
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
