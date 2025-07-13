package deviceMonitorDaoImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorDao"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type IDeviceDataReportDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewIDeviceDataReportDaoImpl(db *gorm.DB) deviceMonitorDao.IDeviceDataReportDao {
	return &IDeviceDataReportDaoImpl{
		db:        db,
		tableName: "sys_iot_db_dev_map",
	}
}

func (i IDeviceDataReportDaoImpl) DevList(c *gin.Context) ([]deviceMonitorModel.SysIotDbDevMapData, error) {
	var list []deviceMonitorModel.SysIotDbDevMapData
	ret := i.db.Table(i.tableName).Where("state = ?", commonStatus.NORMAL).Limit(500).Find(&list)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return list, nil
}

func (i IDeviceDataReportDaoImpl) GetDevList(c *gin.Context, dev []string) ([]deviceMonitorModel.SysIotDbDevMapData, error) {
	var list []deviceMonitorModel.SysIotDbDevMapData
	ret := i.db.Table(i.tableName).Where("device in (?)", dev).Where("state = ?", commonStatus.NORMAL).Limit(500).Find(&list)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return list, nil
}

func (i IDeviceDataReportDaoImpl) Save(c *gin.Context, data *deviceMonitorModel.SysIotDbDevMap) error {
	if data == nil {
		return nil
	}

	var value *deviceMonitorModel.SysIotDbDevMapData
	ret := i.db.Table(i.tableName).Where("device = ?", data.Device).Where("state = ?", commonStatus.NORMAL).First(&value)
	if ret.Error != nil && !errors.Is(ret.Error, gorm.ErrRecordNotFound) {
		zap.L().Error("save device map error", zap.Error(ret.Error))
		return ret.Error
	}
	if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
		value = nil
	}
	if value == nil {
		data.ID = snowflake.GenID()
		data.SetCreateBy(baizeContext.GetUserId(c))
		ret = i.db.Table(i.tableName).Create(data)
		if ret.Error != nil {
			return ret.Error
		}
	} else {
		data.SetUpdateBy(baizeContext.GetUserId(c))
		ret = i.db.Table(i.tableName).Where("device = ?", data.Device).Updates(data)
		return ret.Error
	}
	return nil
}

func (i IDeviceDataReportDaoImpl) Remove(c *gin.Context, dev string) error {
	ret := i.db.Table(i.tableName).Where("device = ?", dev).Update("state", commonStatus.DELETE)
	return ret.Error
}
