package masterdaoimpl

import (
	"nova-factory-server/app/business/erp/master/masterdao"

	"gorm.io/gorm"
)

type StockDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewStockDao(db *gorm.DB) masterdao.IStockDao {
	return &StockDaoImpl{db: db, tableName: "erp_stock"}
}
