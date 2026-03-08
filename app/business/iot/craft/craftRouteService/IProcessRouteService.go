package craftRouteService

import (
	"nova-factory-server/app/business/iot/craft/craftRouteModels"

	"github.com/gin-gonic/gin"
)

type IProcessRouteService interface {
	Add(c *gin.Context, request *craftRouteModels.SysProRouteProcessSetRequest) (*craftRouteModels.SysProRouteProcess, error)

	Update(c *gin.Context, request *craftRouteModels.SysProRouteProcessSetRequest) (*craftRouteModels.SysProRouteProcess, error)

	List(c *gin.Context, req *craftRouteModels.SysProRouteProcessListReq) (*craftRouteModels.SysProRouteProcessList, error)

	Remove(c *gin.Context, ids []string) error
}
