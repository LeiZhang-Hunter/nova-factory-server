package craftRouteDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

type ISysCraftRouteConfigDao interface {
	GetById(routeId uint64) (*craftRouteModels.SysCraftRouteConfig, error)
	Save(c *gin.Context, routeId uint64, topo *craftRouteModels.ProcessTopo) (*craftRouteModels.SysCraftRouteConfig, error)
}
