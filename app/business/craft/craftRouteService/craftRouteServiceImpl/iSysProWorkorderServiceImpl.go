package craftRouteServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
)

type ISysProWorkorderServiceImpl struct {
	dao craftRouteDao.ISysProWorkorderDao
}

func NewISysProWorkorderServiceImpl(dao craftRouteDao.ISysProWorkorderDao) craftRouteService.ISysProWorkorderService {
	return &ISysProWorkorderServiceImpl{
		dao: dao,
	}
}

func (i *ISysProWorkorderServiceImpl) Add(c *gin.Context, info *craftRouteModels.SysSetProWorkorder) (*craftRouteModels.SysProWorkorder, error) {
	return i.dao.Add(c, info)
}
func (i *ISysProWorkorderServiceImpl) Update(c *gin.Context, info *craftRouteModels.SysSetProWorkorder) (*craftRouteModels.SysProWorkorder, error) {
	return i.dao.Update(c, info)
}
func (i *ISysProWorkorderServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
func (i *ISysProWorkorderServiceImpl) List(c *gin.Context, req *craftRouteModels.SysProWorkorderReq) (*craftRouteModels.SysProWorkorderList, error) {
	return i.dao.List(c, req)
}
