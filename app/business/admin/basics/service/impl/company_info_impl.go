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

func (s *CompanyInfoService) SelectCompanyInfo(c *gin.Context) *models.CompanyInfoVo {
	return s.cd.SelectCompanyInfo(c)
}

func (s *CompanyInfoService) SaveCompanyInfo(c *gin.Context, company *models.CompanyInfoVo) {
	company.CompanyName = strings.TrimSpace(company.CompanyName)

	if !s.cd.ExistsCompanyInfo(c) {
		s.cd.InsertCompanyInfo(c, company)
		return
	}

	s.cd.UpdateCompanyInfo(c, company)
}
