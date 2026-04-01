package craftroutedao

import (
	"nova-factory-server/app/business/iot/craft/craftroutemodels"

	"github.com/gin-gonic/gin"
)

type IProcessContextDao interface {
	Add(c *gin.Context, processContext *craftroutemodels.SysProSetProcessContent) (*craftroutemodels.SysProProcessContent, error)
	Update(c *gin.Context, processContext *craftroutemodels.SysProSetProcessContent) (*craftroutemodels.SysProProcessContent, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, process *craftroutemodels.SysProProcessContextListReq) (*craftroutemodels.SysProProcessContextListData, error)
	GetByProcessIds(c *gin.Context, ids []int64) ([]*craftroutemodels.SysProProcessContent, error)
}
