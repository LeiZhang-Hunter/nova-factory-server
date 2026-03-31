package devicedaoImpl

import (
	"errors"
	"nova-factory-server/app/business/iot/asset/device/devicedao"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ISysModbusDeviceConfigDataDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewISysModbusDeviceConfigDataDaoImp(ms *gorm.DB) devicedao.ISysModbusDeviceConfigDataDao {
	return &ISysModbusDeviceConfigDataDaoImpl{
		db:        ms,
		tableName: "sys_modbus_device_config_data",
	}
}

func (i *ISysModbusDeviceConfigDataDaoImpl) Add(c *gin.Context, data *devicemodels.SysModbusDeviceConfigData) (*devicemodels.SysModbusDeviceConfigData, error) {
	ret := i.db.Table(i.tableName).Create(data)
	return data, ret.Error
}

func (i *ISysModbusDeviceConfigDataDaoImpl) Update(c *gin.Context, data *devicemodels.SysModbusDeviceConfigData) (*devicemodels.SysModbusDeviceConfigData, error) {
	ret := i.db.Table(i.tableName).Debug().Where("device_config_id = ?", data.DeviceConfigID).UpdateColumns(data)
	return data, ret.Error
}

func (i *ISysModbusDeviceConfigDataDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := i.db.Table(i.tableName).Where("device_config_id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (i *ISysModbusDeviceConfigDataDaoImpl) List(c *gin.Context, req *devicemodels.SysModbusDeviceConfigDataListReq) (*devicemodels.SysModbusDeviceConfigDataListData, error) {
	db := i.db.Table(i.tableName)
	if req == nil {
		req = &devicemodels.SysModbusDeviceConfigDataListReq{}
	}
	if req.TemplateID != 0 {
		db = db.Where("template_id = ?", req.TemplateID)
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
		return &devicemodels.SysModbusDeviceConfigDataListData{
			Rows:  make([]*devicemodels.SetSysModbusDeviceConfigDataReq, 0),
			Total: 0,
		}, ret.Error
	}
	var dto []*devicemodels.SysModbusDeviceConfigData
	ret = db.Offset(offset).Limit(size).Order("create_time desc").Find(&dto)
	if ret.Error != nil {
		return &devicemodels.SysModbusDeviceConfigDataListData{
			Rows:  make([]*devicemodels.SetSysModbusDeviceConfigDataReq, 0),
			Total: 0,
		}, ret.Error
	}
	var rows []*devicemodels.SetSysModbusDeviceConfigDataReq = make([]*devicemodels.SetSysModbusDeviceConfigDataReq, 0)
	for _, data := range dto {
		rows = append(rows, devicemodels.ToSetSysModbusDeviceConfigDataReq(data))
	}
	return &devicemodels.SysModbusDeviceConfigDataListData{
		Rows:  rows,
		Total: total,
	}, nil
}

func (i *ISysModbusDeviceConfigDataDaoImpl) GetByTemplateIds(c *gin.Context, ids []uint64) ([]*devicemodels.SysModbusDeviceConfigData, error) {
	if ids == nil || len(ids) == 0 {
		return nil, errors.New("ids is null")
	}
	var dto []*devicemodels.SysModbusDeviceConfigData
	ret := i.db.Table(i.tableName).Where("template_id in (?)", ids).Where("state = ?", commonStatus.NORMAL).Find(&dto)
	if ret.Error != nil {
		return nil, ret.Error
	}

	return dto, nil
}

func (i *ISysModbusDeviceConfigDataDaoImpl) GetById(c *gin.Context, id uint64) (*devicemodels.SysModbusDeviceConfigData, error) {

	var dto *devicemodels.SysModbusDeviceConfigData
	ret := i.db.Table(i.tableName).Where("device_config_id = ?", id).Where("state = ?", commonStatus.NORMAL).First(&dto)
	if ret.Error != nil {
		return nil, ret.Error
	}

	return dto, nil
}

func (i *ISysModbusDeviceConfigDataDaoImpl) GetByIds(c *gin.Context, id []uint64) ([]*devicemodels.SysModbusDeviceConfigData, error) {

	var dto []*devicemodels.SysModbusDeviceConfigData
	ret := i.db.Table(i.tableName).Where("device_config_id in (?)", id).Where("state = ?", commonStatus.NORMAL).Find(&dto)
	if ret.Error != nil {
		return nil, ret.Error
	}

	return dto, nil
}
