package dao

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IApiShopCompanyInfoDao provides mini app company information reads.
type IApiShopCompanyInfoDao interface {
	Get(c *gin.Context) (*models.CompanyInfo, error)
}
