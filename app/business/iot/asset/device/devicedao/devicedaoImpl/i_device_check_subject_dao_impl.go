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

type IDeviceCheckSubjectDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewIDeviceCheckSubjectDaoImpl(db *gorm.DB) devicedao.IDeviceCheckSubjectDao {
	return &IDeviceCheckSubjectDaoImpl{
		db:    db,
		table: "sys_device_check_subject",
	}
}

func (i *IDeviceCheckSubjectDaoImpl) Set(c *gin.Context, data *devicemodels.SysDeviceCheckSubjectVO) (*devicemodels.SysDeviceCheckSubject, error) {
	if data == nil {
		return nil, errors.New("data should not be nil")
	}
	value := devicemodels.ToSysDeviceCheckSubject(data)
	if data.RecordID != 0 {
		value.SetUpdateBy(baizeContext.GetUserId(c))
		ret := i.db.Table(i.table).Where("plan_id = ?", data.PlanID).UpdateColumns(value)
		return value, ret.Error
	} else {
		value.RecordID = snowflake.GenID()
		value.DeptID = baizeContext.GetDeptId(c)
		value.SetCreateBy(baizeContext.GetUserId(c))
		ret := i.db.Table(i.table).Create(value)
		return value, ret.Error
	}
}

func (i *IDeviceCheckSubjectDaoImpl) List(c *gin.Context, req *devicemodels.SysDeviceCheckSubjectReq) (*devicemodels.SysDeviceCheckSubjectList, error) {
	db := i.db.Table(i.table)
	if req == nil {
		req = &devicemodels.SysDeviceCheckSubjectReq{}
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
		return &devicemodels.SysDeviceCheckSubjectList{
			Rows:  make([]*devicemodels.SysDeviceCheckSubject, 0),
			Total: 0,
		}, ret.Error
	}
	var dto []*devicemodels.SysDeviceCheckSubject
	ret = db.Offset(offset).Limit(size).Order("create_time desc").Find(&dto)
	if ret.Error != nil {
		return &devicemodels.SysDeviceCheckSubjectList{
			Rows:  make([]*devicemodels.SysDeviceCheckSubject, 0),
			Total: 0,
		}, ret.Error
	}
	return &devicemodels.SysDeviceCheckSubjectList{
		Rows:  dto,
		Total: total,
	}, nil
}

func (i *IDeviceCheckSubjectDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := i.db.Table(i.table).Where("id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}
