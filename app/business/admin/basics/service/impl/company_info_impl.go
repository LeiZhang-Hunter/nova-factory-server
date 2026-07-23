package impl

import (
	"nova-factory-server/app/business/admin/basics/dao"
	"nova-factory-server/app/business/admin/basics/models"
	"nova-factory-server/app/business/admin/basics/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type CompanyInfoService struct {
	cd dao.ICompanyInfoDao
}

func NewCompanyInfoService(cd dao.ICompanyInfoDao) service.ICompanyInfoService {
	return &CompanyInfoService{cd: cd}
}

func (s *CompanyInfoService) SelectCompanyInfo(c *gin.Context) (*models.CompanyInfoVo, error) {
	return s.cd.SelectCompanyInfo(c)
}

func (s *CompanyInfoService) SaveCompanyInfo(c *gin.Context, company *models.CompanyInfoVo) error {
	company.CompanyName = strings.TrimSpace(company.CompanyName)

	exist, err := s.cd.ExistsCompanyInfo(c)
	if err != nil {
		return err
	}
	if !exist {
		return s.cd.InsertCompanyInfo(c, company)
	}

	return s.cd.UpdateCompanyInfo(c, company)
}
