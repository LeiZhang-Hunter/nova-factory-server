package shopdaoimpl

import (
	"errors"
	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/constant/commonStatus"

	"gorm.io/gorm"
)

// IShopOrderDaoImpl 商城订单主表 DAO 实现。
//
// 该 DAO 只负责 shop_order 单表：
// - 开启事务；
// - 按 tid 查询主表；
// - 创建主表；
// - 按 id 更新主表。
//
// 不在这里编排订单明细和订单账户，也不在这里判断订单状态是否允许覆盖。
// 这些属于订单同步业务流程，由 service 层统一处理。
type IShopOrderDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewIShopOrderDaoImpl 创建商城订单主表 DAO。
func NewIShopOrderDaoImpl(db *gorm.DB) shopdao.IShopOrderDao {
	return &IShopOrderDaoImpl{
		db:        db,
		tableName: "shop_order",
	}
}

// Transaction 开启订单同步事务。
//
// service 层会在该事务中组合调用主表 DAO、明细 DAO、账户 DAO。
// fn 返回 error 时 GORM 自动回滚；fn 返回 nil 时提交。
func (i *IShopOrderDaoImpl) Transaction(fn func(tx *gorm.DB) error) error {
	if i.db == nil {
		return errors.New("shop order dao db is nil")
	}
	return i.db.Transaction(fn)
}

// GetByTid 根据 tid 查询有效订单主表。
//
// 找不到记录时返回 nil, nil，方便 service 层判断新增或更新。
func (i *IShopOrderDaoImpl) GetByTid(tx *gorm.DB, tid string) (*shopmodels.Order, error) {
	var item shopmodels.Order
	if err := tx.Table(i.tableName).
		Where("tid = ?", tid).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// Create 创建订单主表。
//
// 只写 shop_order，明细和账户由 service 层在同一事务中调用对应 DAO 写入。
func (i *IShopOrderDaoImpl) Create(tx *gorm.DB, order *shopmodels.Order) error {
	return tx.Table(i.tableName).Create(order).Error
}

// UpdateByID 按 id 更新订单主表。
//
// updates 由 service 层构建。这里不做业务判断，只执行主表更新。
func (i *IShopOrderDaoImpl) UpdateByID(tx *gorm.DB, id uint64, updates map[string]any) error {
	return tx.Table(i.tableName).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		Updates(updates).Error
}
