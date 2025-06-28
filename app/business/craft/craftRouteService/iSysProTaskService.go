package craftRouteService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

type ISysProTaskService interface {
	Add(ctx *gin.Context, req *craftRouteModels.SysSetProTask) (*craftRouteModels.SysProTask, error)
	Update(ctx *gin.Context, req *craftRouteModels.SysSetProTask) (*craftRouteModels.SysProTask, error)
	Remove(ctx *gin.Context, ids []string) error
	List(ctx *gin.Context, req *craftRouteModels.SysProTaskReq) (*craftRouteModels.SysProTaskList, error)
}
