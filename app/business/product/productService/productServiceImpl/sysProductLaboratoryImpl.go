package productServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/product/productDao"
	"nova-factory-server/app/business/product/productModels"
	"nova-factory-server/app/business/product/productService"
)

type SysProductLaboratoryService struct {
	dao productDao.ISysProductLaboratoryDao
}

func NewSysProductLaboratoryService(dao productDao.ISysProductLaboratoryDao) productService.ISysProductLaboratoryService {
	return &SysProductLaboratoryService{dao: dao}
}

func (s *SysProductLaboratoryService) SelectLaboratoryList(c *gin.Context, dql *productModels.SysProductLaboratoryDQL) (list *productModels.SysProductLaboratoryList, err error) {
	return s.dao.SelectLaboratoryList(c, dql)
}

func (s *SysProductLaboratoryService) SelectLaboratoryById(c *gin.Context, id int64) (*productModels.SysProductLaboratoryVo, error) {
	return s.dao.SelectLaboratoryById(c, id)
}

func (s *SysProductLaboratoryService) Set(c *gin.Context, data *productModels.SysProductLaboratoryVo) (*productModels.SysProductLaboratory, error) {
	return s.dao.Set(c, data)
}

func (s *SysProductLaboratoryService) DeleteLaboratoryByIds(c *gin.Context, ids []int64) error {
	err := s.dao.DeleteLaboratoryByIds(c, ids)
	return err
}
