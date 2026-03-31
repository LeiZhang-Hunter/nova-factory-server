package craftrouteservice

import (
	"nova-factory-server/app/business/iot/craft/craftroutemodels"

	"github.com/gin-gonic/gin"
)

type ICraftProcessService interface {
	Add(c *gin.Context, process *craftroutemodels.SysProSetProcessReq) (*craftroutemodels.SysProProcess, error)
	Update(c *gin.Context, process *craftroutemodels.SysProSetProcessReq) (*craftroutemodels.SysProProcess, error)
	Remove(c *gin.Context, processIds []int64) error
	List(c *gin.Context, process *craftroutemodels.SysProProcessListReq) (*craftroutemodels.SysProProcessListData, error)
}
