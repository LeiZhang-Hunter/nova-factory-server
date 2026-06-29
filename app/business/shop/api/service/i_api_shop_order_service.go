package service

import (
	"nova-factory-server/app/business/shop/api/models"
	shopordermodels "nova-factory-server/app/business/shop/order/models"

	"github.com/gin-gonic/gin"
)

// IApiShopOrderService   订单服务接口
type IApiShopOrderService interface {
	Confirm(c *gin.Context, userID int64, req *models.OrderConfirmReq) (*models.OrderConfirmResp, error)
	Create(c *gin.Context, userID int64, req *models.OrderCreateReq) (*shopordermodels.Order, error)
	GetByID(c *gin.Context, id int64) (*models.ApiOrderVO, error)
	List(c *gin.Context, userID int64, query *models.OrderQuery) (*models.OrderListData, error)
	UpdateStatus(c *gin.Context, userID int64, req *models.OrderStatusReq) error
	Pay(c *gin.Context, userID int64, id int64, payChannel int) (*models.OrderPayResp, error)
	//HandleWechatNotify(event.ZOrderStatusSyncReqEvent) error
	Cancel(c *gin.Context, userID int64, id int64, reason string) error
	ConfirmReceive(c *gin.Context, userID int64, id int64) error
	GetStatistics(c *gin.Context, userID int64) (*models.OrderStatistics, error)
	BatchUpdateStatus(c *gin.Context, userID int64, req []models.OrderStatus) error
}
