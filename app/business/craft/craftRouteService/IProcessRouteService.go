package craftRouteService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

type IProcessRouteService interface {
	Add(c *gin.Context, request *craftRouteModels.SysProRouteProcessSetRequest) (*craftRouteModels.SysProRouteProcess, error)

	Update(c *gin.Context, request *craftRouteModels.SysProRouteProcessSetRequest) (*craftRouteModels.SysProRouteProcess, error)

	List(c *gin.Context, req *craftRouteModels.SysProRouteProcessListReq) (*craftRouteModels.SysProRouteProcessList, error)

	Remove(c *gin.Context, ids []string) error
}
