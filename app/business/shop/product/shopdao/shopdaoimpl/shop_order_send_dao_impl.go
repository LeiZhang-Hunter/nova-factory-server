package shopdaoimpl

import (
	"errors"

	"gorm.io/gorm"
	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
)

// ShopOrderSendDaoImpl 订单发货主表 DAO 实现。
//
// 该 DAO 只负责 shop_order_send 单表：
// - 开启事务；
// - 创建主表记录。
//
// 不在这里编排发货明细，明细写入由 service 层在同一事务中组合调用。
type ShopOrderSendDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewShopOrderSendDaoImpl 创建发货主表 DAO。
func NewShopOrderSendDaoImpl(db *gorm.DB) shopdao.IShopOrderSendDao {
	return &ShopOrderSendDaoImpl{
		db:        db,
		tableName: "shop_order_send",
	}
}

// Transaction 开启发货事务。
//
// service 层会在该事务中组合调用发货主表 DAO、发货明细 DAO。
func (d *ShopOrderSendDaoImpl) Transaction(fn func(tx *gorm.DB) error) error {
	if d.db == nil {
		return errors.New("shop order send dao db is nil")
	}
	return d.db.Transaction(fn)
}

// Create 新增发货主表记录。
func (d *ShopOrderSendDaoImpl) Create(tx *gorm.DB, send *shopmodels.OrderSend) error {
	return tx.Table(d.tableName).Create(send).Error
}
