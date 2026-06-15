package shopdao

import (
	"gorm.io/gorm"
	"nova-factory-server/app/business/shop/product/shopmodels"
)

// IShopOrderSendDao 订单发货主表数据访问接口。
//
// DAO 层只负责 shop_order_send 单表操作和事务能力。
// 业务编排（发货主表和明细表的组合写入）由 service 层负责。
type IShopOrderSendDao interface {
	// Transaction 开启发货事务。
	//
	// service 层会在该事务中组合调用发货主表 DAO、发货明细 DAO。
	Transaction(fn func(tx *gorm.DB) error) error

	// Create 新增发货主表记录。
	Create(tx *gorm.DB, send *shopmodels.OrderSend) error
}
