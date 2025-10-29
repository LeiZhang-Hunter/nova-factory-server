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

type IDeviceSubjectDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewIDeviceSubjectDaoImpl(db *gorm.DB) deviceDao.IDeviceSubjectDao {
	return &IDeviceSubjectDaoImpl{
		db:    db,
		table: "sys_device_subject",
	}
}

func (i *IDeviceSubjectDaoImpl) Set(c *gin.Context, data *deviceModels.SysDeviceSubjectVO) (*deviceModels.SysDeviceSubject, error) {
	if data == nil {
		return nil, errors.New("data should not be nil")
	}
	value := deviceModels.ToSysDeviceSubject(data)
	if data.ID != 0 {
		value.SetUpdateBy(baizeContext.GetUserId(c))
		ret := i.db.Table(i.table).Debug().UpdateColumns(value)
		return value, ret.Error
	} else {
		value.ID = snowflake.GenID()
		value.DeptID = baizeContext.GetDeptId(c)
		value.SetCreateBy(baizeContext.GetUserId(c))
		ret := i.db.Table(i.table).Where("id = ?", data.ID).Create(value)
		return value, ret.Error
	}
}

func (i *IDeviceSubjectDaoImpl) List(c *gin.Context, req *deviceModels.SysDeviceSubjectReq) (*deviceModels.SysDeviceSubjectListData, error) {
	db := i.db.Table(i.table)
	if req == nil {
		req = &deviceModels.SysDeviceSubjectReq{}
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
		return &deviceModels.SysDeviceSubjectListData{
			Rows:  make([]*deviceModels.SysDeviceSubject, 0),
			Total: 0,
		}, ret.Error
	}
	var dto []*deviceModels.SysDeviceSubject
	ret = db.Offset(offset).Limit(size).Order("create_time desc").Find(&dto)
	if ret.Error != nil {
		return &deviceModels.SysDeviceSubjectListData{
			Rows:  make([]*deviceModels.SysDeviceSubject, 0),
			Total: 0,
		}, ret.Error
	}
	return &deviceModels.SysDeviceSubjectListData{
		Rows:  dto,
		Total: total,
	}, nil
}

func (i *IDeviceSubjectDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := i.db.Table(i.table).Where("id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}
