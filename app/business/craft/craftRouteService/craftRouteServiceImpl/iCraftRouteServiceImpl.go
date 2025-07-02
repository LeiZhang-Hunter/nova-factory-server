package craftRouteServiceImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
)

type CraftRouteServiceImpl struct {
	dao             craftRouteDao.ICraftRouteDao
	processDao      craftRouteDao.IProcessDao
	routeProcessDao craftRouteDao.IRouteProcessDao
	contextDao      craftRouteDao.IProcessContextDao
	bomDao          craftRouteDao.ISysProRouteProductBomDao
	productDao      craftRouteDao.ISysProRouteProductDao
	routeConfigDao  craftRouteDao.ISysCraftRouteConfigDao
}

func NewCraftRouteServiceImpl(dao craftRouteDao.ICraftRouteDao,
	processDao craftRouteDao.IProcessDao,
	routeProcessDao craftRouteDao.IRouteProcessDao,
	contextDao craftRouteDao.IProcessContextDao,
	bomDao craftRouteDao.ISysProRouteProductBomDao,
	productDao craftRouteDao.ISysProRouteProductDao,
	routeConfigDao craftRouteDao.ISysCraftRouteConfigDao) craftRouteService.ICraftRouteService {
	return &CraftRouteServiceImpl{
		dao:             dao,
		processDao:      processDao,
		routeProcessDao: routeProcessDao,
		contextDao:      contextDao,
		bomDao:          bomDao,
		productDao:      productDao,
		routeConfigDao:  routeConfigDao,
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

func (craft *CraftRouteServiceImpl) DetailCraftRoute(c *gin.Context, req *craftRouteModels.SysCraftRouteDetailRequest) (*craftRouteModels.SysCraftRouteConfig, error) {
	// 读取详情
	info, err := craft.routeConfigDao.GetById(uint64(req.RouteID))
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (craft *CraftRouteServiceImpl) SaveCraftRoute(c *gin.Context, topo *craftRouteModels.ProcessTopo) (*craftRouteModels.SysCraftRouteConfig, error) {
	if topo == nil {
		return nil, errors.New("参数错误")
	}
	if topo.Route == nil {
		return nil, errors.New("工艺路线参数错误")
	}
	info, err := craft.dao.GetById(c, topo.Route.RouteID)
	if err != nil {
		zap.L().Error("读取工艺路线失败", zap.Error(err))
		return nil, err
	}
	if info == nil {
		return nil, errors.New("工艺路线不存在")
	}

	data, err := craft.routeConfigDao.Save(c, uint64(topo.Route.RouteID), topo)
	if err != nil {
		return nil, err
	}
	return data, err
}
