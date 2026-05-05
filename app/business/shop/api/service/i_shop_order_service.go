package service

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IShopOrderService 订单服务接口
type IShopOrderService interface {
	Create(c *gin.Context, username string, req *models.OrderSetReq) (*models.Order, error)
	GetByID(c *gin.Context, id int64) (*models.OrderVO, error)
	List(c *gin.Context, username string, query *models.OrderQuery) (*models.OrderListData, error)
	UpdateStatus(c *gin.Context, username string, req *models.OrderStatusReq) error
	Cancel(c *gin.Context, username string, id int64, reason string) error
	ConfirmReceive(c *gin.Context, username string, id int64) error
	GetStatistics(c *gin.Context, username string) (*models.OrderStatistics, error)
}
