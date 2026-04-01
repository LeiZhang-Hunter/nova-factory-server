package deviceservice

import (
	"nova-factory-server/app/business/iot/asset/device/devicemodels"

	"github.com/gin-gonic/gin"
)

type IDeviceService interface {
	InsertDevice(c *gin.Context, job *devicemodels.DeviceInfo) (*devicemodels.DeviceVO, error)
	UpdateDevice(c *gin.Context, job *devicemodels.DeviceInfo) (*devicemodels.DeviceVO, error)
	SelectDeviceList(c *gin.Context, req *devicemodels.DeviceListReq) (*devicemodels.DeviceInfoListValue, error)
	DeleteByDeviceIds(c *gin.Context, ids []int64) error
	GetById(c *gin.Context, id int64) (*devicemodels.DeviceVO, error)
	GetMetricByTag(c *gin.Context, req *devicemodels.DeviceTagListReq) (*devicemodels.DeviceMetricInfoListValue, error)
	// GetDeviceInfoByIds 通过id 列表读取设备状态，包含设备在线状态以及异常状态
	GetDeviceInfoByIds(c *gin.Context, ids []int64) ([]*devicemodels.DeviceVO, error)
	// StatCount 统计设备，在线、不在线、异常
	StatCount(c *gin.Context) (*devicemodels.DeviceStatData, error)
}
