package deviceservice

import (
	"nova-factory-server/app/business/iot/asset/device/devicemodels"

	"github.com/gin-gonic/gin"
)

type IDeviceGroupService interface {
	InsertDeviceGroup(c *gin.Context, group *devicemodels.DeviceGroup) (*devicemodels.DeviceGroupVO, error)
	UpdateDeviceGroup(c *gin.Context, group *devicemodels.DeviceGroup) (*devicemodels.DeviceGroupVO, error)
	SelectDeviceGroupList(c *gin.Context, req *devicemodels.DeviceGroupDQL) (*devicemodels.DeviceGroupListData, error)
	DeleteByGroupIds(c *gin.Context, ids []int64) error
}
