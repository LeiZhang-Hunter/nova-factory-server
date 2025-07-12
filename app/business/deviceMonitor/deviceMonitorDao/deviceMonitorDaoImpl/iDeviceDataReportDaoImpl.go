package deviceMonitorDaoImpl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorDao"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/constant/commonStatus"
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
