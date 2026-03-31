package devicedaoImpl

import (
	"errors"
	"nova-factory-server/app/business/iot/asset/device/devicedao"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IDeviceCheckPlanDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewIDeviceCheckPlanDaoImpl(db *gorm.DB) devicedao.IDeviceCheckPlanDao {
	return &IDeviceCheckPlanDaoImpl{
		db:    db,
		table: "sys_device_check_plan",
	}
}

func (i *IDeviceCheckPlanDaoImpl) Set(c *gin.Context, data *devicemodels.SysDeviceCheckPlanVO) (*devicemodels.SysDeviceCheckPlan, error) {
	if data == nil {
		return nil, errors.New("data should not be nil")
	}
	value := devicemodels.ToSysDeviceCheckPlan(data)
	if data.PlanID != 0 {
		value.SetUpdateBy(baizeContext.GetUserId(c))
		ret := i.db.Table(i.table).Where("plan_id = ?", data.PlanID).UpdateColumns(value)
		return value, ret.Error
	} else {
		value.PlanID = snowflake.GenID()
		value.DeptID = baizeContext.GetDeptId(c)
		value.SetCreateBy(baizeContext.GetUserId(c))
		ret := i.db.Table(i.table).Create(value)
		return value, ret.Error
	}
}

func (i *IDeviceCheckPlanDaoImpl) List(c *gin.Context, req *devicemodels.SysDeviceCheckPlanReq) (*devicemodels.SysDeviceCheckPlanList, error) {
	db := i.db.Table(i.table)
	if req == nil {
		req = &devicemodels.SysDeviceCheckPlanReq{}
	}
	if req != nil && req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
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

	db = baizeContext.GetGormDataScope(c, db)
	db = db.Where("state = ?", 0)

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &devicemodels.SysDeviceCheckPlanList{
			Rows:  make([]*devicemodels.SysDeviceCheckPlan, 0),
			Total: 0,
		}, ret.Error
	}
	var dto []*devicemodels.SysDeviceCheckPlan
	ret = db.Offset(offset).Limit(size).Order("create_time desc").Find(&dto)
	if ret.Error != nil {
		return &devicemodels.SysDeviceCheckPlanList{
			Rows:  make([]*devicemodels.SysDeviceCheckPlan, 0),
			Total: 0,
		}, ret.Error
	}
	return &devicemodels.SysDeviceCheckPlanList{
		Rows:  dto,
		Total: total,
	}, nil
}

func (i *IDeviceCheckPlanDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := i.db.Table(i.table).Where("id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}
