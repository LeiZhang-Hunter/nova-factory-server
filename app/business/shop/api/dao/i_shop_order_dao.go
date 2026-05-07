package dao

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IApiShopOrderDao  订单数据访问接口
type IApiShopOrderDao interface {
	Create(c *gin.Context, order *models.Order) (*models.Order, error)
	GetByID(c *gin.Context, id int64) (*models.Order, error)
	GetByOrderNo(c *gin.Context, orderNo string) (*models.Order, error)
	List(c *gin.Context, query *models.OrderQuery) (*models.OrderListData, error)
	UpdateStatus(c *gin.Context, id int64, status int32, version int32) (int64, error)
	Cancel(c *gin.Context, id int64, reason string, version int32) (int64, error)
	GetStatistics(c *gin.Context, userID int64) (*models.OrderStatistics, error)
}

// IApiShopOrderItemDao  订单商品明细数据访问接口
type IApiShopOrderItemDao interface {
	BatchCreate(c *gin.Context, items []*models.OrderItem) error
	GetByOrderID(c *gin.Context, orderID int64) ([]*models.OrderItem, error)
	GetByOrderNo(c *gin.Context, orderNo string) ([]*models.OrderItem, error)
}
