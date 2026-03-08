package craftRouteDao

import (
	"nova-factory-server/app/business/iot/craft/craftRouteModels"

	"github.com/gin-gonic/gin"
)

type IProcessContextDao interface {
	Add(c *gin.Context, processContext *craftRouteModels.SysProSetProcessContent) (*craftRouteModels.SysProProcessContent, error)
	Update(c *gin.Context, processContext *craftRouteModels.SysProSetProcessContent) (*craftRouteModels.SysProProcessContent, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, process *craftRouteModels.SysProProcessContextListReq) (*craftRouteModels.SysProProcessContextListData, error)
	GetByProcessIds(c *gin.Context, ids []int64) ([]*craftRouteModels.SysProProcessContent, error)
}
