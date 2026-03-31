package craftRouteDao

import (
	"nova-factory-server/app/business/iot/craft/craftRouteModels"
	"nova-factory-server/app/business/iot/craft/craftRouteModels/api/v1"

	"github.com/gin-gonic/gin"
)

type ISysCraftRouteConfigDao interface {
	GetById(routeId uint64) (*craftRouteModels.SysCraftRouteConfig, error)
	GetConfigByIds(routeIds []int64) ([]*v1.Router, error)
	Save(c *gin.Context, routeId uint64, topo *craftRouteModels.ProcessTopo, config []byte) (*craftRouteModels.SysCraftRouteConfig, error)
}
