package craftrouteservice

import (
	"nova-factory-server/app/business/iot/craft/craftroutemodels"

	"github.com/gin-gonic/gin"
)

type ICraftRouteService interface {
	AddCraftRoute(c *gin.Context, route *craftroutemodels.SysCraftRouteRequest) (*craftroutemodels.SysCraftRoute, error)
	UpdateCraftRoute(c *gin.Context, route *craftroutemodels.SysCraftRouteRequest) (*craftroutemodels.SysCraftRoute, error)
	RemoveCraftRoute(c *gin.Context, ids []int64) error
	SelectCraftRoute(c *gin.Context, req *craftroutemodels.SysCraftRouteListReq) (*craftroutemodels.SysCraftRouteListData, error)
	DetailCraftRoute(c *gin.Context, req *craftroutemodels.SysCraftRouteDetailRequest) (*craftroutemodels.SysCraftRouteConfig, error)
	SaveCraftRoute(c *gin.Context, topo *craftroutemodels.ProcessTopo) (*craftroutemodels.SysCraftRouteConfig, error)
	// Count 统计调度任务数量
	Count(c *gin.Context) (int64, error)
}
