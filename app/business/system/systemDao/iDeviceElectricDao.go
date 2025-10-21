package systemDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/system/systemModels"
)

type IDeviceElectricDao interface {
	Set(c *gin.Context, setting *systemModels.SysDeviceElectricSettingVO) (*systemModels.SysDeviceElectricSetting, error)
	List(c *gin.Context, req *systemModels.SysDeviceElectricSettingDQL) (*systemModels.SysDeviceElectricSettingData, error)
	Remove(c *gin.Context, ids []string) error
	GetByDeviceId(c *gin.Context, id int64) (*systemModels.SysDeviceElectricSetting, error)
	GetByNoDeviceId(c *gin.Context, id int64, deviceId int64) (*systemModels.SysDeviceElectricSetting, error)
	GetByIds(c *gin.Context, ids []string) ([]*systemModels.SysDeviceElectricSetting, error)
	All(c *gin.Context) ([]*systemModels.SysDeviceElectricSetting, error)
}
