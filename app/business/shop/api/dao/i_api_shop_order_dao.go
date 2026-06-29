package dao

import (
	shopordermodels "nova-factory-server/app/business/shop/order/models"
	shopusermodels "nova-factory-server/app/business/shop/user/models"
	"time"

	models "nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IApiShopOrderDao interface {
	// ListShopOrders 分页查询当前商城用户的 shop_order 订单。
	ListShopOrders(c *gin.Context, shopUser *shopusermodels.User, query *models.OrderQuery) (*models.OrderListData, error)
	// UpdateERPOrderStatus 更新当前商城用户的 shop_order 订单状态。
	UpdateERPOrderStatus(c *gin.Context, id int64, shopUser *shopusermodels.User, status int32) (int64, error)
	// CancelERPOrder 将当前商城用户的 shop_order 订单标记为已取消。
	CancelERPOrder(c *gin.Context, id int64, shopUser *shopusermodels.User, reason string) (int64, error)
	// GetERPOrderStatistics 统计当前商城用户的 shop_order 订单状态数量。
	GetERPOrderStatistics(c *gin.Context, shopUser *shopusermodels.User) (*models.OrderStatistics, error)
	// MarkOrderPaidWithTx 在事务内锁定并标记 shop_order 为已支付。
	MarkOrderPaidWithTx(c *gin.Context, tx *gorm.DB, id uint64, payTime *time.Time, transactionId, notifyRaw, mchId, appid, payerOpenid string) error

	// GetByTid 按订单编号查询 shop 订单详情。
	GetByTid(c *gin.Context, tid string) (*shopordermodels.Order, error)

	GetByID(c *gin.Context, id uint64) (*shopordermodels.Order, error)
	// GetByIDs 根据id批量读取订单
	GetByIDs(c *gin.Context, ids []int64) ([]*shopordermodels.Order, error)

	// BatchUpdateERPOrderStatus 批量更新当前商城用户的 shop_order 订单状态。
	BatchUpdateERPOrderStatus(c *gin.Context, ids []int64, shopUser *shopusermodels.User, status int32) (int64, error)

	//
}
