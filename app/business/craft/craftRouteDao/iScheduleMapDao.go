package craftRouteDao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

type IScheduleMapDao interface {
	GetByScheduleIds(c *gin.Context, ids []int64) ([]*craftRouteModels.SysProductScheduleMap, error)
	GetSpecialSchedule(c *gin.Context, beginTime int64) ([]*craftRouteModels.SysProductScheduleMap, error)
	Set(c *gin.Context, tx *gorm.DB, data *craftRouteModels.SetSysProductSchedule) error
	Remove(c *gin.Context, ids []string) error
	GetByScheduleId(c *gin.Context, id int64) ([]*craftRouteModels.SysProductScheduleMap, error)
	GetSpecialScheduleByNow(c *gin.Context, gatewayId int64) ([]*craftRouteModels.SysProductScheduleMap, error)
	GetNormalByTime(c *gin.Context, gatewayId int64) ([]*craftRouteModels.SysProductScheduleMap, error)
}
