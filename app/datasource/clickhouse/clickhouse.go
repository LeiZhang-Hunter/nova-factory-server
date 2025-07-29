package clickhouse

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type ClickHouse struct {
	db *gorm.DB
}

func NewClickHouse() (*ClickHouse, error) {
	return nil, nil
	dsn := viper.GetString("clickhouse.link")
	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("connect clickhouse fail", zap.Error(err))
		return nil, nil
	}
	return &ClickHouse{db: db}, nil
}

func (c *ClickHouse) DB() *gorm.DB {
	return c.db
}
