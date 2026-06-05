package impl

import (
	"context"
	"nova-factory-server/app/business/admin/basics/dao"
	"nova-factory-server/app/business/admin/basics/models"

	"gorm.io/gorm"
)

type companyInfoDao struct {
	db        *gorm.DB
	tableName string
}

func NewCompanyInfoDao(db *gorm.DB) dao.ICompanyInfoDao {
	return &companyInfoDao{
		db:        db,
		tableName: "sys_company_info",
	}
}

func (d *companyInfoDao) SelectCompanyInfo(ctx context.Context) *models.CompanyInfoVo {
	company := new(models.CompanyInfoVo)
	err := d.db.WithContext(ctx).Table(d.tableName).Limit(1).Find(company).Error
	if err != nil {
		panic(err)
	}
	return company
}

func (d *companyInfoDao) ExistsCompanyInfo(ctx context.Context) bool {
	var count int64
	err := d.db.WithContext(ctx).Table(d.tableName).Count(&count).Error
	if err != nil {
		panic(err)
	}
	return count > 0
}

func (d *companyInfoDao) InsertCompanyInfo(ctx context.Context, company *models.CompanyInfoVo) {
	err := d.db.WithContext(ctx).Table(d.tableName).Create(company).Error
	if err != nil {
		panic(err)
	}
}

func (d *companyInfoDao) UpdateCompanyInfo(ctx context.Context, company *models.CompanyInfoVo) {
	updates := map[string]interface{}{
		"company_name":        company.CompanyName,
		"company_detail":      company.CompanyDetail,
		"contact_phone":       company.ContactPhone,
		"contact_person":      company.ContactPerson,
		"email":               company.Email,
		"address":             company.Address,
		"logo_url":            company.LogoUrl,
		"business_license_no": company.BusinessLicenseNo,
		"remark":              company.Remark,
	}
	err := d.db.WithContext(ctx).
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Table(d.tableName).
		Updates(updates).Error
	if err != nil {
		panic(err)
	}
}
