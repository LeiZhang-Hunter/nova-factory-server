package deviceDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceModels"
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
