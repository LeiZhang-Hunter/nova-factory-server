package dao

import (
	"context"
	"nova-factory-server/app/business/admin/basics/models"
)

type ICompanyInfoDao interface {
	SelectCompanyInfo(ctx context.Context) *models.CompanyInfoVo
	ExistsCompanyInfo(ctx context.Context) bool
	InsertCompanyInfo(ctx context.Context, company *models.CompanyInfoVo)
	UpdateCompanyInfo(ctx context.Context, company *models.CompanyInfoVo)
}
