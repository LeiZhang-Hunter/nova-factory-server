package service

import (
	"nova-factory-server/app/business/admin/basics/models"

	"github.com/gin-gonic/gin"
)

type ICompanyInfoService interface {
	SelectCompanyInfo(c *gin.Context) *models.CompanyInfoVo
	SaveCompanyInfo(c *gin.Context, company *models.CompanyInfoVo)
}
