package craftRouteDao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

type IScheduleDao interface {
	GetDailySchedule(c *gin.Context) ([]*craftRouteModels.SysProductSchedule, error)
	List(c *gin.Context, req *craftRouteModels.SysProductScheduleListReq) (*craftRouteModels.SysProductScheduleListData, error)
	Set(c *gin.Context, tx *gorm.DB, data *craftRouteModels.SetSysProductSchedule) (*craftRouteModels.SysProductSchedule, error)
	Remove(c *gin.Context, ids []string) error
	GetById(c *gin.Context, id int64) (*craftRouteModels.SysProductSchedule, error)
}
