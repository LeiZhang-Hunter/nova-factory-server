package craftroutedao

import (
	"nova-factory-server/app/business/iot/craft/craftroutemodels"
	"nova-factory-server/app/business/iot/craft/craftroutemodels/api/v1"

	"github.com/gin-gonic/gin"
)

type ISysCraftRouteConfigDao interface {
	GetById(routeId uint64) (*craftroutemodels.SysCraftRouteConfig, error)
	GetConfigByIds(routeIds []int64) ([]*v1.Router, error)
	Save(c *gin.Context, routeId uint64, topo *craftroutemodels.ProcessTopo, config []byte) (*craftroutemodels.SysCraftRouteConfig, error)
}
