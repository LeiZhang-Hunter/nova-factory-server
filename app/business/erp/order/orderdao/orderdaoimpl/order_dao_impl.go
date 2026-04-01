package orderdaoimpl

import (
	"nova-factory-server/app/business/erp/order/orderdao"

	"gorm.io/gorm"
)

type OrderDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewOrderDao(db *gorm.DB) orderdao.IOrderDao {
	return &OrderDaoImpl{
		db: db,
	}
}
