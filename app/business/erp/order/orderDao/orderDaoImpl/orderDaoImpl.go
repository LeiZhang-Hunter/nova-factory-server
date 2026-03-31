package orderDaoImpl

import (
	"nova-factory-server/app/business/erp/order/orderDao"

	"gorm.io/gorm"
)

type OrderDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewOrderDao(db *gorm.DB) orderDao.IOrderDao {
	return &OrderDaoImpl{
		db: db,
	}
}
