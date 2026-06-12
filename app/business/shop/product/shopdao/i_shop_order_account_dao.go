package shopdao

import (
	"gorm.io/gorm"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"time"
)

// IShopOrderAccountDao 商城订单账户表数据访问接口。
//
// 只负责 shop_order_account 单表操作。调用方必须传入外层事务 tx，
// 保证账户表操作能和订单主表、明细表一起提交或回滚。
type IShopOrderAccountDao interface {
	BatchCreate(tx *gorm.DB, orderID uint64, order *shopmodels.Order, now *time.Time) error
	DeleteByOrderID(tx *gorm.DB, orderID uint64) error
}
