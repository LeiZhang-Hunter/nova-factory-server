package deviceService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceModels"
)

type IDeviceGroupService interface {
	InsertDeviceGroup(c *gin.Context, group *deviceModels.DeviceGroup) (*deviceModels.DeviceGroupVO, error)
	UpdateDeviceGroup(c *gin.Context, group *deviceModels.DeviceGroup) (*deviceModels.DeviceGroupVO, error)
	SelectDeviceGroupList(c *gin.Context, req *deviceModels.DeviceGroupDQL) (*deviceModels.DeviceGroupListData, error)
	DeleteByGroupIds(c *gin.Context, ids []int64) error
}
