package dashboardDaoImpl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/dashboard/dashboardDao"
	"nova-factory-server/app/business/dashboard/dashboardModels"
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
		value.ID = snowflake.GenID()
		value.DeptID = baizeContext.GetDeptId(c)
		value.SetCreateBy(baizeContext.GetUserId(c))
		ret := d.db.Table(d.table).Create(value)
		return value, ret.Error
	}
	value.SetUpdateBy(baizeContext.GetUserId(c))
	ret := d.db.Table(d.table).Where("id = ?", data.ID).Updates(value)
	return value, ret.Error
}
func (d *DashboardDataDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := d.db.Table(d.table).Where("id in (?)", ids).Delete(&dashboardModels.SetSysDashboardData{})
	return ret.Error
}
