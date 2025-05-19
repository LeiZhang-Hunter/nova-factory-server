package craftRouteServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
)

type ISysProRouteProductBomServiceImpl struct {
	dao craftRouteDao.ISysProRouteProductBomDao
}

func NewISysProRouteProductBomServiceImpl(dao craftRouteDao.ISysProRouteProductBomDao) craftRouteService.ISysProRouteProductBomService {
	return &ISysProRouteProductBomServiceImpl{
		dao: dao,
	}
}

func (i *ISysProRouteProductBomServiceImpl) Add(c *gin.Context, info *craftRouteModels.SysSetProRouteProductBom) (*craftRouteModels.SysProRouteProductBom, error) {
	return i.dao.Add(c, info)
}

func (i *ISysProRouteProductBomServiceImpl) Update(c *gin.Context, info *craftRouteModels.SysSetProRouteProductBom) (*craftRouteModels.SysProRouteProductBom, error) {
	return i.dao.Update(c, info)
}
func (i *ISysProRouteProductBomServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
func (i *ISysProRouteProductBomServiceImpl) List(c *gin.Context, req *craftRouteModels.SysProRouteProductBomReq) (*craftRouteModels.SysProRouteProductBomList, error) {
	return i.dao.List(c, req)
}
