package dashboardDaoImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/dashboard/dashboardDao"
	"nova-factory-server/app/business/dashboard/dashboardModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type DashboardDataDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewDashboardDataDaoImpl(db *gorm.DB) dashboardDao.DashboardDataDao {
	return &DashboardDataDaoImpl{
		db:    db,
		table: "sys_dashboard_data",
	}
}

func (d *DashboardDataDaoImpl) Set(c *gin.Context, data *dashboardModels.SetSysDashboardData) (*dashboardModels.SysDashboardData, error) {
	value := dashboardModels.ToSysDashboardData(data)
	if data.ID == 0 {
		var info *dashboardModels.SysDashboardData
		ret := d.db.Table(d.table).Where("datashboard_id = ?", data.DatashboardID).Where("state = ?", commonStatus.NORMAL).First(&info)
		if ret.Error != nil && !errors.Is(gorm.ErrRecordNotFound, ret.Error) {
			zap.L().Error("get dashboard error", zap.Error(ret.Error))
			return nil, ret.Error
		}
		if info == nil || info.ID == 0 {
			value.ID = snowflake.GenID()
			value.DeptID = baizeContext.GetDeptId(c)
			value.SetCreateBy(baizeContext.GetUserId(c))
			ret = d.db.Table(d.table).Create(value)
			return value, ret.Error
		}
	}
	value.SetUpdateBy(baizeContext.GetUserId(c))
	ret := d.db.Table(d.table).Where("datashboard_id = ?", data.DatashboardID).Updates(value)
	return value, ret.Error
}
func (d *DashboardDataDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := d.db.Table(d.table).Where("id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (d *DashboardDataDaoImpl) Info(c *gin.Context, dashboardId int64) (*dashboardModels.SysDashboardData, error) {
	var info *dashboardModels.SysDashboardData
	ret := d.db.Table(d.table).Where("datashboard_id = ?", dashboardId).Where("state = ?", commonStatus.NORMAL).First(&info)
	if ret.Error != nil && !errors.Is(gorm.ErrRecordNotFound, ret.Error) {
		return nil, ret.Error
	}
	if info.ID == 0 {
		return nil, nil
	}
	return info, nil
}
