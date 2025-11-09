package craftRouteService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
	v1 "nova-factory-server/app/business/craft/craftRouteModels/api/v1"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
)

type IScheduleService interface {
	GetMonthSchedule(c *gin.Context, req *craftRouteModels.SysProductScheduleReq) ([]*craftRouteModels.ScheduleStatusData, error)
	List(c *gin.Context, req *craftRouteModels.SysProductScheduleListReq) (*craftRouteModels.SysProductScheduleListData, error)
	Set(c *gin.Context, data *craftRouteModels.SetSysProductSchedule) error
	Remove(c *gin.Context, ids []string) error
	Detail(c *gin.Context, id int64) (*craftRouteModels.DetailSysProductData, error)
	Schedule(ctx *gin.Context, req *craftRouteModels.ScheduleReq, gatewayInfo *daemonizeModels.SysIotAgent) ([]*v1.Router, error)
}
