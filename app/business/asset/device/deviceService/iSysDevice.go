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
}
