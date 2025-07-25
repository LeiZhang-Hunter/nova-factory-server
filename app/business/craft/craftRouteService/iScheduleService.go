package craftRouteService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

type IScheduleService interface {
	GetMonthSchedule(c *gin.Context, req *craftRouteModels.SysProductScheduleReq) ([]*craftRouteModels.ScheduleStatusData, error)
	List(c *gin.Context, req *craftRouteModels.SysProductScheduleListReq) (*craftRouteModels.SysProductScheduleListData, error)
	Set(c *gin.Context, data *craftRouteModels.SetSysProductSchedule) error
	Remove(c *gin.Context, ids []string) error
	Detail(c *gin.Context, id int64) (*craftRouteModels.DetailSysProductData, error)
}
