package craftRouteService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

type ICraftRouteService interface {
	AddCraftRoute(c *gin.Context, route *craftRouteModels.SysCraftRouteRequest) (*craftRouteModels.SysCraftRoute, error)
	UpdateCraftRoute(c *gin.Context, route *craftRouteModels.SysCraftRouteRequest) (*craftRouteModels.SysCraftRoute, error)
	RemoveCraftRoute(c *gin.Context, ids []int64) error
	SelectCraftRoute(c *gin.Context, req *craftRouteModels.SysCraftRouteListReq) (*craftRouteModels.SysCraftRouteListData, error)
	DetailCraftRoute(c *gin.Context, req *craftRouteModels.SysCraftRouteDetailRequest) (*craftRouteModels.SysCraftRouteConfig, error)
	SaveCraftRoute(c *gin.Context, topo *craftRouteModels.ProcessTopo) (*craftRouteModels.SysCraftRouteConfig, error)
	// Count 统计调度任务数量
	Count(c *gin.Context) (int64, error)
}
