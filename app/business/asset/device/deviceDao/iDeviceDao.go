package deviceDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceModels"
)

type IDeviceDao interface {
	InsertDevice(c *gin.Context, device *deviceModels.DeviceInfo) (*deviceModels.DeviceVO, error)
	UpdateDevice(c *gin.Context, device *deviceModels.DeviceInfo) (*deviceModels.DeviceVO, error)
	GetDeviceGroupByName(c *gin.Context, name string) (*deviceModels.DeviceVO, error)
	GetNoExitIdDeviceGroupByName(c *gin.Context, name string, id uint64) (*deviceModels.DeviceVO, error)
	SelectDeviceList(c *gin.Context, req *deviceModels.DeviceListReq) (*deviceModels.DeviceInfoListData, error)
	DeleteByDeviceIds(c *gin.Context, ids []int64) error
}
