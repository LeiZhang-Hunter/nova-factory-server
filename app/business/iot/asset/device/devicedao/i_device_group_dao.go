package devicedao

import (
	"nova-factory-server/app/business/iot/asset/device/devicemodels"

	"github.com/gin-gonic/gin"
)

type IDeviceGroupDao interface {
	InsertDeviceGroup(c *gin.Context, device *devicemodels.DeviceGroup) (*devicemodels.DeviceGroupVO, error)
	UpdateDeviceGroup(c *gin.Context, device *devicemodels.DeviceGroup) (*devicemodels.DeviceGroupVO, error)
	GetDeviceGroupByName(c *gin.Context, name string) (*devicemodels.DeviceGroupVO, error)
	GetNoExitIdDeviceGroupByName(c *gin.Context, name string, id uint64) (*devicemodels.DeviceGroupVO, error)
	SelectDeviceGroupList(c *gin.Context, req *devicemodels.DeviceGroupDQL) (*devicemodels.DeviceGroupListData, error)
	GetDeviceGroupByIds(c *gin.Context, ids []uint64) ([]*devicemodels.DeviceGroupVO, error)
	DeleteByGroupIds(c *gin.Context, ids []int64) error
}
