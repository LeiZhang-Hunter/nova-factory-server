package deviceDaoImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type IDeviceCheckMachineryDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewIDeviceCheckMachineryDaoImpl(db *gorm.DB) deviceDao.IDeviceCheckMachineryDao {
	return &IDeviceCheckMachineryDaoImpl{
		db:    db,
		table: "sys_device_check_machinery",
	}
}

func (i *IDeviceCheckMachineryDaoImpl) Set(c *gin.Context, data *deviceModels.SysDeviceCheckMachineryVO) (*deviceModels.SysDeviceCheckMachinery, error) {
	if data == nil {
		return nil, errors.New("data should not be nil")
	}
	value := deviceModels.ToSysDeviceCheckMachinery(data)
	if data.RecordID != 0 {
		value.SetUpdateBy(baizeContext.GetUserId(c))
		ret := i.db.Table(i.table).Where("record_id = ?", data.RecordID).UpdateColumns(value)
		return value, ret.Error
	} else {
		value.PlanID = snowflake.GenID()
		value.DeptID = baizeContext.GetDeptId(c)
		value.SetCreateBy(baizeContext.GetUserId(c))
		ret := i.db.Table(i.table).Create(value)
		return value, ret.Error
	}
}

func (i *IDeviceCheckMachineryDaoImpl) List(c *gin.Context, req *deviceModels.SysDeviceCheckMachineryReq) (*deviceModels.SysDeviceCheckMachineryList, error) {
	db := i.db.Table(i.table)
	if req == nil {
		req = &deviceModels.SysDeviceCheckMachineryReq{}
	}
	if req != nil && req.PlanID != 0 {
		db = db.Where("plan_id = ?", req.PlanID)
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
		return &deviceModels.SysDeviceCheckMachineryList{
			Rows:  make([]*deviceModels.SysDeviceCheckMachinery, 0),
			Total: 0,
		}, ret.Error
	}
	var dto []*deviceModels.SysDeviceCheckMachinery
	ret = db.Offset(offset).Limit(size).Order("create_time desc").Find(&dto)
	if ret.Error != nil {
		return &deviceModels.SysDeviceCheckMachineryList{
			Rows:  make([]*deviceModels.SysDeviceCheckMachinery, 0),
			Total: 0,
		}, ret.Error
	}
	return &deviceModels.SysDeviceCheckMachineryList{
		Rows:  dto,
		Total: total,
	}, nil
}

func (i *IDeviceCheckMachineryDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := i.db.Table(i.table).Where("id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}
