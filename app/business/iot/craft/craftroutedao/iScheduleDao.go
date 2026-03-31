package craftroutedao

import (
	"nova-factory-server/app/business/iot/craft/craftroutemodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IScheduleDao interface {
	GetDailySchedule(c *gin.Context) ([]*craftroutemodels.SysProductSchedule, error)
	List(c *gin.Context, req *craftroutemodels.SysProductScheduleListReq) (*craftroutemodels.SysProductScheduleListData, error)
	Set(c *gin.Context, tx *gorm.DB, data *craftroutemodels.SetSysProductSchedule) (*craftroutemodels.SysProductSchedule, error)
	Remove(c *gin.Context, ids []string) error
	GetById(c *gin.Context, id int64) (*craftroutemodels.SysProductSchedule, error)
}
