package craftroutedao

import (
	craftRouteModels2 "nova-factory-server/app/business/iot/craft/craftroutemodels"
	"nova-factory-server/app/business/iot/daemonize/daemonizemodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IScheduleMapDao interface {
	GetByScheduleIds(c *gin.Context, ids []int64) ([]*craftRouteModels2.SysProductScheduleMap, error)
	GetSpecialSchedule(c *gin.Context, beginTime int64) ([]*craftRouteModels2.SysProductScheduleMap, error)
	Set(c *gin.Context, tx *gorm.DB, data *craftRouteModels2.SetSysProductSchedule) error
	Remove(c *gin.Context, ids []string) error
	GetByScheduleId(c *gin.Context, id int64) ([]*craftRouteModels2.SysProductScheduleMap, error)
	GetSpecialScheduleByNow(c *gin.Context, gatewayId int64, gatewayInfo *daemonizemodels.SysIotAgent) ([]*craftRouteModels2.SysProductScheduleMap, error)
	GetNormalByTime(c *gin.Context, gatewayId int64) ([]*craftRouteModels2.SysProductScheduleMap, error)
}
