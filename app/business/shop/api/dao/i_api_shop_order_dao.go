package dao

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/erp/sale/salemodels"
	models "nova-factory-server/app/business/shop/api/models"
)

type IApiShopOrderDao interface {
	// Set 新增或修改 ERP 订单及其子表。
	Set(c *gin.Context, req *salemodels.OrderSet) (*salemodels.Order, error)
	// SetWithTx 新增或修改 ERP 订单及其子表（带事务）。
	SetWithTx(c *gin.Context, tx *gorm.DB, req *salemodels.OrderSet) (*salemodels.Order, error)
	// GetByID 查询 ERP 订单详情。
	GetByID(c *gin.Context, id uint64) (*salemodels.Order, error)
	// List 分页查询 ERP 订单。
	List(c *gin.Context, req *salemodels.OrderQuery) (*salemodels.OrderListData, error)
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
}
