package craftroutedao

import (
	"nova-factory-server/app/business/iot/craft/craftroutemodels"

	"github.com/gin-gonic/gin"
)

type IProcessDao interface {
	List(c *gin.Context, process *craftroutemodels.SysProProcessListReq) (*craftroutemodels.SysProProcessListData, error)
	Add(c *gin.Context, process *craftroutemodels.SysProSetProcessReq) (*craftroutemodels.SysProProcess, error)
	Update(c *gin.Context, process *craftroutemodels.SysProSetProcessReq) (*craftroutemodels.SysProProcess, error)
	Remove(c *gin.Context, processIds []int64) error
	GetById(c *gin.Context, id int64) (*craftroutemodels.SysProProcess, error)
	GetByIds(c *gin.Context, ids []int64) ([]*craftroutemodels.SysProProcess, error)
}
