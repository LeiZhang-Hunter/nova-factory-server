package deviceService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceModels"
)

type IDeviceService interface {
	InsertDevice(c *gin.Context, job *deviceModels.DeviceInfo) (*deviceModels.DeviceVO, error)
	UpdateDevice(c *gin.Context, job *deviceModels.DeviceInfo) (*deviceModels.DeviceVO, error)
	SelectDeviceList(c *gin.Context, req *deviceModels.DeviceListReq) (*deviceModels.DeviceInfoListValue, error)
	DeleteByDeviceIds(c *gin.Context, ids []int64) error
	GetById(c *gin.Context, id int64) (*deviceModels.DeviceVO, error)
	GetMetricByTag(c *gin.Context, req *deviceModels.DeviceTagListReq) (*deviceModels.DeviceMetricInfoListValue, error)
	// GetDeviceInfoByIds 通过id 列表读取设备状态，包含设备在线状态以及异常状态
	GetDeviceInfoByIds(c *gin.Context, ids []int64) ([]*deviceModels.DeviceVO, error)
}
