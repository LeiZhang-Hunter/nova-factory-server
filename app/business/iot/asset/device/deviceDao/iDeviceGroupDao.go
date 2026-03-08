package deviceDao

import (
	"nova-factory-server/app/business/iot/asset/device/deviceModels"

	"github.com/gin-gonic/gin"
)

type IDeviceGroupDao interface {
	InsertDeviceGroup(c *gin.Context, device *deviceModels.DeviceGroup) (*deviceModels.DeviceGroupVO, error)
	UpdateDeviceGroup(c *gin.Context, device *deviceModels.DeviceGroup) (*deviceModels.DeviceGroupVO, error)
	GetDeviceGroupByName(c *gin.Context, name string) (*deviceModels.DeviceGroupVO, error)
	GetNoExitIdDeviceGroupByName(c *gin.Context, name string, id uint64) (*deviceModels.DeviceGroupVO, error)
	SelectDeviceGroupList(c *gin.Context, req *deviceModels.DeviceGroupDQL) (*deviceModels.DeviceGroupListData, error)
	GetDeviceGroupByIds(c *gin.Context, ids []uint64) ([]*deviceModels.DeviceGroupVO, error)
	DeleteByGroupIds(c *gin.Context, ids []int64) error
}
