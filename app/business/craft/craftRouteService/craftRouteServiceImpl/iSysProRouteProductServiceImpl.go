package craftRouteServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
)

type SysProRouteProductServiceImpl struct {
	dao craftRouteDao.ISysProRouteProductDao
}

func NewSysProRouteProductServiceImpl(dao craftRouteDao.ISysProRouteProductDao) craftRouteService.ISysProRouteProductService {
	return &SysProRouteProductServiceImpl{
		dao: dao,
	}
}

func (s *SysProRouteProductServiceImpl) Add(c *gin.Context, data *craftRouteModels.SysProRouteSetProduct) (*craftRouteModels.SysProRouteProduct, error) {
	return s.dao.Add(c, data)
}

func (s *SysProRouteProductServiceImpl) Update(c *gin.Context, data *craftRouteModels.SysProRouteSetProduct) (*craftRouteModels.SysProRouteProduct, error) {
	return s.dao.Update(c, data)
}

func (s *SysProRouteProductServiceImpl) List(c *gin.Context, req *craftRouteModels.SysProRouteProductReq) (*craftRouteModels.SysProRouteProductList, error) {
	return s.dao.List(c, req)
}

func (s *SysProRouteProductServiceImpl) Remove(c *gin.Context, ids []string) error {
	return s.dao.Remove(c, ids)
}
