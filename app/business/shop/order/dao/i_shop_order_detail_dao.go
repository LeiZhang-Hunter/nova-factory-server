package dao

import (
	"nova-factory-server/app/business/shop/order/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IOrderDetailDao shop订单明细数据访问接口。
//
// 该接口只提供明细表的单表操作能力。订单同步时由 service 层控制事务，
// 并组合订单主表 DAO、明细 DAO、账户 DAO 完成完整同步。
type IOrderDetailDao interface {
	// BatchCreate 批量创建保存接口传入的订单明细。
	BatchCreate(tx *gorm.DB, c *gin.Context, orderID uint64, tid string, details []*models.OrderDetailSet) error
	// BatchCreateByOrder 批量创建事件同步转换后的订单明细。
	BatchCreateByOrder(tx *gorm.DB, orderID uint64, order *models.Order, now *time.Time) error
	// DeleteByOrderID 按订单 ID 删除明细记录。
	DeleteByOrderID(tx *gorm.DB, orderID uint64) error
	// DeleteByOrderIDs 按订单 ID 集合删除明细记录。
	DeleteByOrderIDs(tx *gorm.DB, orderIDs []uint64) error
	// DeleteByTidAndOIDs 按订单编号和明细 OID 删除旧明细记录。
	DeleteByTidAndOIDs(tx *gorm.DB, tid string, details []*models.OrderDetailSet) error
	// ListByOrderIDs 按订单 ID 集合查询明细记录。
	ListByOrderIDs(c *gin.Context, orderIDs []uint64) ([]*models.OrderDetail, error)
	ListByOrderID(c *gin.Context, orderIDs uint64) ([]*models.OrderDetail, error)

	// ListByTidTx 在事务内按订单编号查询该订单所有有效明细。
	ListByTidTx(tx *gorm.DB, tid string) ([]*models.OrderDetail, error)

	// IncrementShippedQty 原子累加指定明细行的已发货数量。
	// 使用 shipped_qty = shipped_qty + qty，MySQL UPDATE 自带行锁保证并发安全。
	// 返回受影响行数，0 表示 oid 不存在或 state 非 NORMAL。
	IncrementShippedQty(tx *gorm.DB, orderID uint64, oid string, qty float64) error
}
