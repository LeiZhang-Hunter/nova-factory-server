package craftRouteDao

import (
	"nova-factory-server/app/business/iot/craft/craftRouteModels"

	"github.com/gin-gonic/gin"
)

type IRouteProcessDao interface {
	Add(c *gin.Context, data *craftRouteModels.SysProRouteProcess) (*craftRouteModels.SysProRouteProcess, error)
	Update(c *gin.Context, data *craftRouteModels.SysProRouteProcess) (*craftRouteModels.SysProRouteProcess, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *craftRouteModels.SysProRouteProcessListReq) (*craftRouteModels.SysProRouteProcessList, error)
	GetByRouteId(c *gin.Context, routeId int64) ([]*craftRouteModels.SysProRouteProcess, error)
}
