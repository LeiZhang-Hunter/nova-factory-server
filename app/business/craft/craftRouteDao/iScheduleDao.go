package craftRouteDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

type IScheduleDao interface {
	GetDailySchedule(c *gin.Context) ([]*craftRouteModels.SysProductSchedule, error)
	List(c *gin.Context, req *craftRouteModels.SysProductScheduleListReq) (*craftRouteModels.SysProductScheduleListData, error)
	Set(c *gin.Context, data *craftRouteModels.SetSysProductSchedule) (*craftRouteModels.SysProductSchedule, error)
}
