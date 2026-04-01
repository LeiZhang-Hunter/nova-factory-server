package craftroutedao

import (
	"nova-factory-server/app/business/iot/craft/craftroutemodels"

	"github.com/gin-gonic/gin"
)

type IRouteProcessDao interface {
	Add(c *gin.Context, data *craftroutemodels.SysProRouteProcess) (*craftroutemodels.SysProRouteProcess, error)
	Update(c *gin.Context, data *craftroutemodels.SysProRouteProcess) (*craftroutemodels.SysProRouteProcess, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *craftroutemodels.SysProRouteProcessListReq) (*craftroutemodels.SysProRouteProcessList, error)
	GetByRouteId(c *gin.Context, routeId int64) ([]*craftroutemodels.SysProRouteProcess, error)
}
