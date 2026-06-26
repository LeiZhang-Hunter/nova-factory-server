package dao

import (
	"gorm.io/gorm"
	"nova-factory-server/app/business/shop/order/models"
)

// IOrderShipmentDao 订单发货物流记录数据访问接口。
type IOrderShipmentDao interface {
	// BatchInsert 在事务内批量插入物流记录。
	BatchInsert(tx *gorm.DB, shipments []*models.OrderShipmentSet) error

	// ExistsByOutsidTx 在事务内按物流单号查重，用于幂等判断。
	ExistsByOutsidTx(tx *gorm.DB, outsid string) (bool, error)

	// ListByOrderIDTx 在事务内按订单 ID 查询所有物流记录。
	ListByOrderIDTx(tx *gorm.DB, orderID uint64) ([]*models.OrderShipment, error)
}
