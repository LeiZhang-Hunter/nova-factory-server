package craftRouteServiceImpl

import (
	"github.com/gin-gonic/gin"
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
}

func NewCraftRouteServiceImpl(dao craftRouteDao.ICraftRouteDao,
	processDao craftRouteDao.IProcessDao,
	routeProcessDao craftRouteDao.IRouteProcessDao,
	contextDao craftRouteDao.IProcessContextDao,
	bomDao craftRouteDao.ISysProRouteProductBomDao,
	productDao craftRouteDao.ISysProRouteProductDao) craftRouteService.ICraftRouteService {
	return &CraftRouteServiceImpl{
		dao:             dao,
		processDao:      processDao,
		routeProcessDao: routeProcessDao,
		contextDao:      contextDao,
		bomDao:          bomDao,
		productDao:      productDao,
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

func (craft *CraftRouteServiceImpl) DetailCraftRoute(c *gin.Context, req *craftRouteModels.SysCraftRouteDetailRequest) (*craftRouteModels.ProcessTopo, error) {
	// 读取详情
	info, err := craft.dao.GetById(c, req.RouteID)
	if err != nil {
		return &craftRouteModels.ProcessTopo{}, err
	}
	if info == nil {
		return &craftRouteModels.ProcessTopo{}, nil
	}

	// 工序组成表
	var processIds []int64 = make([]int64, 0)
	processRouteRelations, err := craft.routeProcessDao.GetByRouteId(c, info.RouteID)
	if err != nil {
		return &craftRouteModels.ProcessTopo{}, err
	}
	if len(processRouteRelations) == 0 {
		return &craftRouteModels.ProcessTopo{}, nil
	}

	for _, value := range processRouteRelations {
		processIds = append(processIds, value.ProcessID)
	}

	// 读取工序列表
	processes, err := craft.processDao.GetByIds(c, processIds)
	if err != nil {
		return &craftRouteModels.ProcessTopo{}, err
	}

	// 读取工序内容列表
	processContexts, err := craft.contextDao.GetByProcessIds(c, processIds)
	if err != nil {
		return &craftRouteModels.ProcessTopo{}, err
	}

	// 读取工序物料列表
	boms, err := craft.bomDao.GetByRouteId(c, info.RouteID)
	if err != nil {
		return &craftRouteModels.ProcessTopo{}, err
	}
	// 读取产品制程
	products, err := craft.productDao.GetByRouteId(c, info.RouteID)
	if err != nil {
		return &craftRouteModels.ProcessTopo{}, err
	}
	topo := craftRouteModels.NewProcessTopo()
	topo.Route = info
	for _, relation := range processRouteRelations {
		// 组装边
		topo.Edges = append(topo.Edges, craftRouteModels.NewProcessTopoEdge(relation))
	}
	// 组装工序内容
	topo.OfProcess(processes, processContexts, boms, products)
	return topo, nil
}
