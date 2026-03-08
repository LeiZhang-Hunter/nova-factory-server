package craftRouteService

import (
	"nova-factory-server/app/business/iot/craft/craftRouteModels"

	"github.com/gin-gonic/gin"
)

type ICraftProcessService interface {
	Add(c *gin.Context, process *craftRouteModels.SysProSetProcessReq) (*craftRouteModels.SysProProcess, error)
	Update(c *gin.Context, process *craftRouteModels.SysProSetProcessReq) (*craftRouteModels.SysProProcess, error)
	Remove(c *gin.Context, processIds []int64) error
	List(c *gin.Context, process *craftRouteModels.SysProProcessListReq) (*craftRouteModels.SysProProcessListData, error)
}
