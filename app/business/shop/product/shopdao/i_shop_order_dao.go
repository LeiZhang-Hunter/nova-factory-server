package shopdao

import (
	"gorm.io/gorm"
	"nova-factory-server/app/business/shop/product/shopmodels"
)

// IShopOrderDao 商城订单主表数据访问接口。
//
// DAO 层只负责 shop_order 单表操作和事务能力，不再编排明细/账户 DAO。
// 增量同步、状态是否允许覆盖、主子表写入顺序等业务逻辑由 service 层负责。
type IShopOrderDao interface {
	// Transaction 开启订单同步事务。
	//
	// service 层会在该事务中组合调用订单主表 DAO、订单明细 DAO、订单账户 DAO。
	// 只要 fn 返回 error，GORM 会回滚整个事务；fn 返回 nil 时才提交。
	Transaction(fn func(tx *gorm.DB) error) error

	// GetByTid 根据订单 tid 查询有效订单。
	//
	// 找不到时返回 nil, nil，方便 service 判断新增还是更新。
	GetByTid(tx *gorm.DB, tid string) (*shopmodels.Order, error)

	// Create 新增订单主表。
	//
	// 只写 shop_order，不写明细和账户。
	Create(tx *gorm.DB, order *shopmodels.Order) error

	// UpdateByID 按主键更新订单主表。
	//
	// updates 由 service 构建，尤其 status 字段必须先经过 service 的状态保护校验。
	UpdateByID(tx *gorm.DB, id uint64, updates map[string]any) error
}
