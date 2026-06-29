package dao

import (
	"nova-factory-server/app/business/shop/order/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IOrderRefundDao 售后单数据访问接口。
type IOrderRefundDao interface {
	Create(c *gin.Context, data *models.OrderRefund) error
	CreateWithTx(tx *gorm.DB, data *models.OrderRefund) error
	GetByID(c *gin.Context, id int64) (*models.OrderRefund, error)
	GetByOutRefundNo(c *gin.Context, outRefundNo string) (*models.OrderRefund, error)
	GetByOrderId(c *gin.Context, orderId int64) (*models.OrderRefund, error)
	UpdateStatus(c *gin.Context, id int64, status int32, updates map[string]any) error
	UpdateByID(c *gin.Context, id int64, updates map[string]any) error
	UpdateStatusWithTx(tx *gorm.DB, id int64, status int32, updates map[string]any) error
	List(c *gin.Context, req *models.RefundQuery) (*models.RefundListData, error)
}
