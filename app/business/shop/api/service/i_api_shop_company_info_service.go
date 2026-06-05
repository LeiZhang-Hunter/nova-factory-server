package service

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IApiShopCompanyInfoService provides public mini app company information.
type IApiShopCompanyInfoService interface {
	Get(c *gin.Context) (*models.CompanyInfoResp, error)
}
