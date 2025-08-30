package craftRouteDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
	v1 "nova-factory-server/app/business/craft/craftRouteModels/api/v1"
)

type ISysCraftRouteConfigDao interface {
	GetById(routeId uint64) (*craftRouteModels.SysCraftRouteConfig, error)
	GetConfigByIds(routeIds []int64) ([]*v1.Router, error)
	Save(c *gin.Context, routeId uint64, topo *craftRouteModels.ProcessTopo, config []byte) (*craftRouteModels.SysCraftRouteConfig, error)
}
