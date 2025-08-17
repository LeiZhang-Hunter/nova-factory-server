package dashboardDaoImpl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/dashboard/dashboardDao"
	"nova-factory-server/app/business/dashboard/dashboardModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type DashboardDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewDashboardDaoImpl(db *gorm.DB) dashboardDao.DashboardDao {
	return &DashboardDaoImpl{
		db:    db,
		table: "sys_dashboard",
	}
}

func (i *DashboardDaoImpl) List(c *gin.Context, req *dashboardModels.SysDashboardReq) (*dashboardModels.SysDashboardList, error) {
	db := i.db.Table(i.table)

	if req != nil && req.Name != "" {
		db = db.Where("name like ?", "%"+req.Name+"%")
	}
	if req != nil && req.Type != "" {
		db = db.Where("type = ?", req.Type)
	}
	size := 0
	if req == nil || req.Size <= 0 {
		size = 20
	} else {
		size = int(req.Size)
	}
	offset := 0
	if req == nil || req.Page <= 0 {
		req.Page = 1
	} else {
		offset = int((req.Page - 1) * req.Size)
	}
	db = db.Where("state", commonStatus.NORMAL)
	db = baizeContext.GetGormDataScope(c, db)
	var dto []*dashboardModels.SysDashboard

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &dashboardModels.SysDashboardList{
			Rows:  make([]*dashboardModels.SysDashboard, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Debug().Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &dashboardModels.SysDashboardList{
			Rows:  make([]*dashboardModels.SysDashboard, 0),
			Total: 0,
		}, ret.Error
	}
	return &dashboardModels.SysDashboardList{
		Rows:  dto,
		Total: total,
	}, nil
}
func (i *DashboardDaoImpl) Set(c *gin.Context, data *dashboardModels.SetSysDashboard) (*dashboardModels.SysDashboard, error) {
	value := dashboardModels.ToSysDashboardReq(data)
	if data.ID == 0 {
		value.ID = snowflake.GenID()
		value.DeptID = baizeContext.GetDeptId(c)
		value.SetCreateBy(baizeContext.GetUserId(c))
		ret := i.db.Table(i.table).Create(value)
		return value, ret.Error
	}
	value.SetUpdateBy(baizeContext.GetUserId(c))
	ret := i.db.Table(i.table).Updates(value)
	return value, ret.Error
}
func (i *DashboardDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := i.db.Table(i.table).Where("id in (?)", ids).Delete(&craftRouteModels.SysProductSchedule{})
	return ret.Error
}
