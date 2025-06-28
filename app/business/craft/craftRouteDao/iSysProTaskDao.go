package craftRouteDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

type ISysProTaskDao interface {
	Add(c *gin.Context, info *craftRouteModels.SysSetProTask) (*craftRouteModels.SysProTask, error)
	Update(c *gin.Context, info *craftRouteModels.SysSetProTask) (*craftRouteModels.SysProTask, error)
	Remove(c *gin.Context, ids []string) error
	List(ctx *gin.Context, req *craftRouteModels.SysProTaskReq) (*craftRouteModels.SysProTaskList, error)
}
