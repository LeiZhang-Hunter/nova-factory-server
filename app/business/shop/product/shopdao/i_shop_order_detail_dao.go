package shopdao

import (
	"gorm.io/gorm"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"time"
)

// IShopOrderDetailDao 商城订单明细表数据访问接口。
//
// 只负责 shop_order_detail 单表操作。调用方必须传入外层事务 tx，
// 保证明细表操作能和订单主表、账户表一起提交或回滚。
type IShopOrderDetailDao interface {
	BatchCreate(tx *gorm.DB, orderID uint64, order *shopmodels.Order, now *time.Time) error
	DeleteByOrderID(tx *gorm.DB, orderID uint64) error
}
