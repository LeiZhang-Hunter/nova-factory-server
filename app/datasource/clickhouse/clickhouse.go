package clickhouse

import (
	"github.com/spf13/viper"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type ClickHouse struct {
	db *gorm.DB
}

func NewClickHouse() (*ClickHouse, error) {
	dsn := viper.GetString("clickhouse.link")
	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &ClickHouse{db: db}, nil
}

func (c *ClickHouse) DB() *gorm.DB {
	return c.db
}
