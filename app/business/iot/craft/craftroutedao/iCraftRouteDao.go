package craftroutedao

import (
	"nova-factory-server/app/business/iot/craft/craftroutemodels"

	"github.com/gin-gonic/gin"
)

type ICraftRouteDao interface {
	// AddCraftRoute 添加工艺
	AddCraftRoute(c *gin.Context, route *craftroutemodels.SysCraftRouteRequest) (*craftroutemodels.SysCraftRoute, error)
	// UpdateCraftRoute 更新工艺
	UpdateCraftRoute(c *gin.Context, route *craftroutemodels.SysCraftRouteRequest) (*craftroutemodels.SysCraftRoute, error)
	// RemoveCraftRoute 移除工艺
	RemoveCraftRoute(c *gin.Context, ids []int64) error
	// SelectCraftRoute 读取工艺列表
	SelectCraftRoute(c *gin.Context, req *craftroutemodels.SysCraftRouteListReq) (*craftroutemodels.SysCraftRouteListData, error)
	GetById(c *gin.Context, id int64) (*craftroutemodels.SysCraftRoute, error)
	GetByIds(c *gin.Context, ids []int64) ([]*craftroutemodels.SysCraftRoute, error)
	// Count 统计调度任务数量
	Count(c *gin.Context) (int64, error)
}
