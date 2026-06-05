package impl

import (
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/utils/fileUtils"

	"github.com/gin-gonic/gin"
)

// IApiShopCompanyInfoServiceImpl provides public mini app company information.
type IApiShopCompanyInfoServiceImpl struct {
	companyDao dao.IApiShopCompanyInfoDao
}

// NewIApiShopCompanyInfoServiceImpl creates the company info service.
func NewIApiShopCompanyInfoServiceImpl(companyDao dao.IApiShopCompanyInfoDao) service.IApiShopCompanyInfoService {
	return &IApiShopCompanyInfoServiceImpl{companyDao: companyDao}
}

// Get returns the public company information payload.
func (s *IApiShopCompanyInfoServiceImpl) Get(c *gin.Context) (*models.CompanyInfoResp, error) {
	company, err := s.companyDao.Get(c)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return &models.CompanyInfoResp{}, nil
	}
	return &models.CompanyInfoResp{
		CompanyName:   company.CompanyName,
		CompanyDetail: company.CompanyDetail,
		ContactPhone:  company.ContactPhone,
		Email:         company.Email,
		Address:       company.Address,
		LogoUrl:       fileUtils.BuildAbsoluteURL(c, company.LogoUrl),
	}, nil
}
