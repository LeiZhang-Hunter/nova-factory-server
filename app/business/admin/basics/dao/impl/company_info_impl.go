package impl

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/admin/basics/dao"
	"nova-factory-server/app/business/admin/basics/models"
	"nova-factory-server/app/utils/fileUtils"

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

func (d *companyInfoDao) SelectCompanyInfo(ctx *gin.Context) (*models.CompanyInfoVo, error) {
	company := new(models.CompanyInfoVo)
	err := d.db.WithContext(ctx).Table(d.tableName).Limit(1).Find(company).Error
	if err != nil {
		return nil, err
	}
	if company != nil && company.LogoUrl != "" {
		company.LogoUrl = fileUtils.BuildAbsoluteURL(ctx, company.LogoUrl)
	}
	return company, nil
}

func (d *companyInfoDao) ExistsCompanyInfo(ctx context.Context) (bool, error) {
	var count int64
	err := d.db.WithContext(ctx).Table(d.tableName).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *companyInfoDao) InsertCompanyInfo(ctx context.Context, company *models.CompanyInfoVo) error {
	if company == nil {
		return nil
	}
	if company.LogoUrl != "" {
		path, err := fileUtils.NormalizeResourcePath(company.LogoUrl)
		if err == nil && path != "" {
			company.LogoUrl = path
		}
	}
	err := d.db.WithContext(ctx).Table(d.tableName).Create(company).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *companyInfoDao) UpdateCompanyInfo(ctx context.Context, company *models.CompanyInfoVo) error {
	if company == nil {
		return nil
	}
	if company.LogoUrl != "" {
		path, err := fileUtils.NormalizeResourcePath(company.LogoUrl)
		if err == nil && path != "" {
			company.LogoUrl = path
		}
	}
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
		return err
	}
	return nil
}
