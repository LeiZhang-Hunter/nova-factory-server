package dao

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/admin/basics/models"
)

type ICompanyInfoDao interface {
	SelectCompanyInfo(ctx *gin.Context) (*models.CompanyInfoVo, error)
	ExistsCompanyInfo(ctx context.Context) (bool, error)
	InsertCompanyInfo(ctx context.Context, company *models.CompanyInfoVo) error
	UpdateCompanyInfo(ctx context.Context, company *models.CompanyInfoVo) error
}
