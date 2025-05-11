package craftRouteServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
)

type CraftRouteServiceImpl struct {
	dao craftRouteDao.ICraftRouteDao
}

func NewCraftRouteServiceImpl(dao craftRouteDao.ICraftRouteDao) craftRouteService.ICraftRouteService {
	return &CraftRouteServiceImpl{
		dao: dao,
	}
}

func (craft *CraftRouteServiceImpl) AddCraftRoute(c *gin.Context, route *craftRouteModels.SysCraftRouteRequest) (*craftRouteModels.SysCraftRoute, error) {
	return craft.dao.AddCraftRoute(c, route)
}

func (craft *CraftRouteServiceImpl) UpdateCraftRoute(c *gin.Context, route *craftRouteModels.SysCraftRouteRequest) (*craftRouteModels.SysCraftRoute, error) {
	return craft.dao.UpdateCraftRoute(c, route)
}

func (craft *CraftRouteServiceImpl) RemoveCraftRoute(c *gin.Context, ids []int64) error {
	return craft.dao.RemoveCraftRoute(c, ids)
}

func (craft *CraftRouteServiceImpl) SelectCraftRoute(c *gin.Context, req *craftRouteModels.SysCraftRouteListReq) (*craftRouteModels.SysCraftRouteListData, error) {
	return craft.dao.SelectCraftRoute(c, req)
}
