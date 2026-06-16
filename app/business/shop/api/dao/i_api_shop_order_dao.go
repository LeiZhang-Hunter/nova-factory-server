package dao

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	models "nova-factory-server/app/business/shop/api/models"
)

type IApiShopOrderDao interface {
	// ListShopOrders 分页查询当前商城用户的 shop_order 订单。
	ListShopOrders(c *gin.Context, shopUser *models.User, query *models.OrderQuery) (*models.OrderListData, error)
	// UpdateERPOrderStatus 更新当前商城用户的 shop_order 订单状态。
	UpdateERPOrderStatus(c *gin.Context, id int64, shopUser *models.User, status int32) (int64, error)
	// CancelERPOrder 将当前商城用户的 shop_order 订单标记为已取消。
	CancelERPOrder(c *gin.Context, id int64, shopUser *models.User, reason string) (int64, error)
	// GetERPOrderStatistics 统计当前商城用户的 shop_order 订单状态数量。
	GetERPOrderStatistics(c *gin.Context, shopUser *models.User) (*models.OrderStatistics, error)
	// MarkOrderPaidWithTx 在事务内锁定并标记 shop_order 为已支付。
	MarkOrderPaidWithTx(c *gin.Context, tx *gorm.DB, id uint64, payTime *time.Time, transactionId, notifyRaw, mchId, appid, payerOpenid string) error
	// DeleteByIDs 删除 ERP 订单。
	DeleteByIDs(c *gin.Context, ids []uint64) error

	// GetByTid 按订单编号查询 shop 订单详情。
	GetByTid(c *gin.Context, tid string) (*models.Order, error)
}
