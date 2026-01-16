package craftRouteDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

type ICraftRouteDao interface {
	// AddCraftRoute 添加工艺
	AddCraftRoute(c *gin.Context, route *craftRouteModels.SysCraftRouteRequest) (*craftRouteModels.SysCraftRoute, error)
	// UpdateCraftRoute 更新工艺
	UpdateCraftRoute(c *gin.Context, route *craftRouteModels.SysCraftRouteRequest) (*craftRouteModels.SysCraftRoute, error)
	// RemoveCraftRoute 移除工艺
	RemoveCraftRoute(c *gin.Context, ids []int64) error
	// SelectCraftRoute 读取工艺列表
	SelectCraftRoute(c *gin.Context, req *craftRouteModels.SysCraftRouteListReq) (*craftRouteModels.SysCraftRouteListData, error)
	GetById(c *gin.Context, id int64) (*craftRouteModels.SysCraftRoute, error)
	GetByIds(c *gin.Context, ids []int64) ([]*craftRouteModels.SysCraftRoute, error)
	// Count 统计调度任务数量
	Count(c *gin.Context) (int64, error)
}
