package craftRouteService

import (
	craftRouteModels2 "nova-factory-server/app/business/iot/craft/craftRouteModels"
	"nova-factory-server/app/business/iot/craft/craftRouteModels/api/v1"
	"nova-factory-server/app/business/iot/daemonize/daemonizeModels"

	"github.com/gin-gonic/gin"
)

type IScheduleService interface {
	GetMonthSchedule(c *gin.Context, req *craftRouteModels2.SysProductScheduleReq) ([]*craftRouteModels2.ScheduleStatusData, error)
	List(c *gin.Context, req *craftRouteModels2.SysProductScheduleListReq) (*craftRouteModels2.SysProductScheduleListData, error)
	Set(c *gin.Context, data *craftRouteModels2.SetSysProductSchedule) error
	Remove(c *gin.Context, ids []string) error
	Detail(c *gin.Context, id int64) (*craftRouteModels2.DetailSysProductData, error)
	Schedule(ctx *gin.Context, req *craftRouteModels2.ScheduleReq, gatewayInfo *daemonizeModels.SysIotAgent) ([]*v1.Router, error)
}
