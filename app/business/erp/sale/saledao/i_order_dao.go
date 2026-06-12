package saledao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/erp/sale/salemodels"
)

// IOrderDao ERP订单数据访问接口。
type IOrderDao interface {
	// Set 新增或修改 ERP 订单及其子表。
	Set(c *gin.Context, req *salemodels.OrderSet) (*salemodels.Order, error)
	// SetWithTx 新增或修改 ERP 订单及其子表（带事务）。
	SetWithTx(c *gin.Context, tx *gorm.DB, req *salemodels.OrderSet) (*salemodels.Order, error)
	// GetByID 查询 ERP 订单详情。
	GetByID(c *gin.Context, id uint64) (*salemodels.Order, error)
	// GetByTid 按订单编号查询 ERP 订单详情。
	GetByTid(c *gin.Context, tid string) (*salemodels.Order, error)
	// List 分页查询 ERP 订单。
	List(c *gin.Context, req *salemodels.OrderQuery) (*salemodels.OrderListData, error)
	// DeleteByIDs 删除 ERP 订单。
	DeleteByIDs(c *gin.Context, ids []uint64) error

	// GetByTidTx 在事务内按订单编号查询 ERP 订单主表。
	GetByTidTx(tx *gorm.DB, tid string) (*salemodels.Order, error)
	// Create 在事务内创建 ERP 订单主表记录。
	Create(tx *gorm.DB, order *salemodels.Order) error
	// UpdateByID 在事务内按 ID 更新 ERP 订单主表记录。
	UpdateByID(tx *gorm.DB, id uint64, updates map[string]any) error

	// Transaction 开启订单同步事务。
	//
	// service 层会在该事务中组合调用订单主表 DAO、订单明细 DAO、订单账户 DAO。
	// 只要 fn 返回 error，GORM 会回滚整个事务；fn 返回 nil 时才提交。
	Transaction(fn func(tx *gorm.DB) error) error
}
