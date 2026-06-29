package service

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IApiShopOrderRefundService  小程序售后接口。
type IApiShopOrderRefundService interface {
	Apply(c *gin.Context, userID int64, req *models.RefundApplyReq) (*models.RefundApplyResp, error)
}
