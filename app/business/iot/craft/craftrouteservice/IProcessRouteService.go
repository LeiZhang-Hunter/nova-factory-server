package craftrouteservice

import (
	"nova-factory-server/app/business/iot/craft/craftroutemodels"

	"github.com/gin-gonic/gin"
)

type IProcessRouteService interface {
	Add(c *gin.Context, request *craftroutemodels.SysProRouteProcessSetRequest) (*craftroutemodels.SysProRouteProcess, error)

	Update(c *gin.Context, request *craftroutemodels.SysProRouteProcessSetRequest) (*craftroutemodels.SysProRouteProcess, error)

	List(c *gin.Context, req *craftroutemodels.SysProRouteProcessListReq) (*craftroutemodels.SysProRouteProcessList, error)

	Remove(c *gin.Context, ids []string) error
}
