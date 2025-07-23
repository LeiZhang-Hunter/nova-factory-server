package craftRouteDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

type IScheduleMapDao interface {
	GetByScheduleIds(c *gin.Context, ids []int64) ([]*craftRouteModels.SysProductScheduleMap, error)
	GetSpecialSchedule(c *gin.Context, beginTime int64) ([]*craftRouteModels.SysProductScheduleMap, error)
	Set(c *gin.Context, data *craftRouteModels.SetSysProductSchedule)
	Remove(c *gin.Context, ids []string) error
	GetByScheduleId(c *gin.Context, id int64) ([]*craftRouteModels.SysProductScheduleMap, error)
}
