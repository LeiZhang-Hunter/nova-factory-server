package service

import (
	"nova-factory-server/app/business/iot/system/models"

	"github.com/gin-gonic/gin"
)

type IDeviceElectricService interface {
	Set(c *gin.Context, setting *models.SysDeviceElectricSettingVO) (*models.SysDeviceElectricSetting, error)
	List(c *gin.Context, req *models.SysDeviceElectricSettingDQL) (*models.SysDeviceElectricSettingData, error)
	Remove(c *gin.Context, ids []string) error
	GetByDeviceId(c *gin.Context, id int64) (*models.SysDeviceElectricSetting, error)
	GetByNoDeviceId(c *gin.Context, id int64, deviceId int64) (*models.SysDeviceElectricSetting, error)
}
