package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/shop/order/models"
)

type IShopOrderShipmentService interface {
	// ListByOrderID 根据id 读取订单列表
	ListByOrderID(ctx *gin.Context, orderID uint64) ([]*models.OrderShipment, error)

	// BatchInsert 批量插入
	BatchInsert(tx *gorm.DB, shipments []*models.OrderShipmentSet) error
}
