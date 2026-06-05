package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IApiShopCompanyInfoDaoImpl reads company information for mini app APIs.
type IApiShopCompanyInfoDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewIApiShopCompanyInfoDaoImpl creates the company info DAO.
func NewIApiShopCompanyInfoDaoImpl(ms *gorm.DB) dao.IApiShopCompanyInfoDao {
	return &IApiShopCompanyInfoDaoImpl{
		db:        ms,
		tableName: "sys_company_info",
	}
}

// Get returns the configured company information.
func (d *IApiShopCompanyInfoDaoImpl) Get(c *gin.Context) (*models.CompanyInfo, error) {
	var company models.CompanyInfo
	if err := d.db.WithContext(c).Table(d.tableName).Take(&company).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}
