package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/shop/order/models"
	"time"
)

// IOrderAccountDao shop订单账户数据访问接口。
//
// 该接口只提供账户表的单表操作能力，事务边界由 service 层或主 DAO 的 SetWithTx 调用方控制。
type IOrderAccountDao interface {
	// BatchCreate 批量创建保存接口传入的订单账户。
	BatchCreate(tx *gorm.DB, c *gin.Context, orderID uint64, tid string, accounts []*models.OrderAccountSet) error
	// BatchCreateByOrder 批量创建事件同步转换后的订单账户。
	BatchCreateByOrder(tx *gorm.DB, orderID uint64, order *models.Order, now *time.Time) error
	// DeleteByOrderID 按订单 ID 删除账户记录。
	DeleteByOrderID(tx *gorm.DB, orderID uint64) error
	// DeleteByOrderIDs 按订单 ID 集合删除账户记录。
	DeleteByOrderIDs(tx *gorm.DB, orderIDs []uint64) error
	// ListByOrderIDs 按订单 ID 集合查询账户记录。
	ListByOrderIDs(c *gin.Context, orderIDs []uint64) ([]*models.OrderAccount, error)
}
