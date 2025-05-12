package craftRouteDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

type IProcessDao interface {
	List(c *gin.Context, process *craftRouteModels.SysProProcessListReq) (*craftRouteModels.SysProProcessListData, error)
	Add(c *gin.Context, process *craftRouteModels.SysProProcess) (*craftRouteModels.SysProProcess, error)
	Update(c *gin.Context, process *craftRouteModels.SysProProcess) (*craftRouteModels.SysProProcess, error)
	Remove(c *gin.Context, processIds []int64) error
}
