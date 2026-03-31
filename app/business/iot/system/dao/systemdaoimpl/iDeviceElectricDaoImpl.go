package systemdaoimpl

import (
	"encoding/json"
	"errors"
	"nova-factory-server/app/business/iot/system/dao"
	"nova-factory-server/app/business/iot/system/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IDeviceElectricDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewIDeviceElectricDaoImpl(db *gorm.DB) dao.IDeviceElectricDao {
	return &IDeviceElectricDaoImpl{
		db:    db,
		table: "sys_device_electric_setting",
	}
}

func (i *IDeviceElectricDaoImpl) Set(c *gin.Context, setting *models.SysDeviceElectricSettingVO) (*models.SysDeviceElectricSetting, error) {
	value := models.ToSysDeviceElectricSetting(setting)
	if value.ID == 0 {
		value.SetCreateBy(baizeContext.GetUserId(c))
		value.ID = snowflake.GenID()
		value.DeptID = baizeContext.GetDeptId(c)
		ret := i.db.Table(i.table).Create(&value)
		return value, ret.Error
	} else {
		value.SetUpdateBy(baizeContext.GetUserId(c))
		ret := i.db.Table(i.table).Debug().Where("id = ?", value.ID).Updates(&value)
		return value, ret.Error
	}
}

func (i *IDeviceElectricDaoImpl) List(c *gin.Context, req *models.SysDeviceElectricSettingDQL) (*models.SysDeviceElectricSettingData, error) {
	db := i.db.Table(i.table)
	if req == nil {
		req = &models.SysDeviceElectricSettingDQL{}
	}
	if len(req.DeviceId) != 0 {
		db = db.Where("device_id = ?", req.DeviceId)
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
	db = db.Where("state = ?", commonStatus.NORMAL)

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &models.SysDeviceElectricSettingData{
			Rows:  make([]*models.SysDeviceElectricSetting, 0),
			Total: 0,
		}, ret.Error
	}
	var dto []*models.SysDeviceElectricSetting
	ret = db.Offset(offset).Limit(size).Order("create_time desc").Find(&dto)
	if ret.Error != nil {
		return &models.SysDeviceElectricSettingData{
			Rows:  make([]*models.SysDeviceElectricSetting, 0),
			Total: 0,
		}, ret.Error
	}
	for k, v := range dto {
		var ex models.Expression
		err := json.Unmarshal([]byte(v.Expression), &ex)
		if err != nil {
			zap.L().Error("json error", zap.Error(err))
			continue
		}
		dto[k].ExpressionData = &ex

	}
	return &models.SysDeviceElectricSettingData{
		Rows:  dto,
		Total: total,
	}, nil
}

func (i *IDeviceElectricDaoImpl) Remove(c *gin.Context, ids []string) error {
	var info models.SysDeviceElectricSetting
	ret := i.db.Table(i.table).Where("id in (?)", ids).Delete(&info)
	return ret.Error
}

func (i *IDeviceElectricDaoImpl) GetByDeviceId(c *gin.Context, id int64) (*models.SysDeviceElectricSetting, error) {
	var info *models.SysDeviceElectricSetting
	ret := i.db.Table(i.table).Where("device_id = ?", id).Where("state = ?", commonStatus.NORMAL).First(&info)
	if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if ret.Error != nil {
		return nil, ret.Error
	}
	return info, nil
}

func (i *IDeviceElectricDaoImpl) GetByNoDeviceId(c *gin.Context, id int64, deviceId int64) (*models.SysDeviceElectricSetting, error) {
	var info *models.SysDeviceElectricSetting
	ret := i.db.Table(i.table).Where("id   != ?", id).Where("device_id = ?", deviceId).Where("state = ?", commonStatus.NORMAL).First(&info)
	if ret.Error != nil {
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, ret.Error
	}
	return info, nil
}

func (i *IDeviceElectricDaoImpl) GetByIds(c *gin.Context, ids []string) ([]*models.SysDeviceElectricSetting, error) {
	var dto []*models.SysDeviceElectricSetting
	ret := i.db.Table(i.table).Debug().Where("id in (?)", ids).Where("state = ?", commonStatus.NORMAL).First(&dto)
	if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if ret.Error != nil {
		return nil, ret.Error
	}
	return dto, nil
}

func (i *IDeviceElectricDaoImpl) All(c *gin.Context) ([]*models.SysDeviceElectricSetting, error) {
	var dto []*models.SysDeviceElectricSetting
	ret := i.db.Table(i.table).Where("state = ?", commonStatus.NORMAL).Find(&dto)
	if ret.Error != nil {
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			return make([]*models.SysDeviceElectricSetting, 0), nil
		}
		return nil, ret.Error
	}
	return dto, nil
}
