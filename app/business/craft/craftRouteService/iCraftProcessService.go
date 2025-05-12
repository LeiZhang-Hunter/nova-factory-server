package craftRouteService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

type ICraftProcessService interface {
	Add(c *gin.Context, process *craftRouteModels.SysProProcess) (*craftRouteModels.SysProProcess, error)
	Update(c *gin.Context, process *craftRouteModels.SysProProcess) (*craftRouteModels.SysProProcess, error)
	Remove(c *gin.Context, processIds []int64) error
	List(c *gin.Context, process *craftRouteModels.SysProProcessListReq) (*craftRouteModels.SysProProcessListData, error)
}
