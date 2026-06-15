package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/shop/order/models"
	"time"
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
}
