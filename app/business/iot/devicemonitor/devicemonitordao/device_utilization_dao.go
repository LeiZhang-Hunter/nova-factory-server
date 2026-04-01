package devicemonitordao

import (
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"

	"github.com/gin-gonic/gin"
)

type DeviceUtilizationDao interface {
	Stat(c *gin.Context, req *devicemonitormodel.DeviceUtilizationReq) ([]*devicemonitormodel.DeviceUtilizationData, error)
	Search(c *gin.Context, req *devicemonitormodel.DeviceUtilizationReq) (*devicemonitormodel.DeviceUtilizationPublicDataList, error)
	SearchV2(c *gin.Context, req *devicemonitormodel.DeviceUtilizationReq) (*devicemonitormodel.DeviceUtilizationPublicDataListV2, error)
	GetDeviceUtilization(c *gin.Context, req *devicemonitormodel.DeviceUtilizationReq) (*devicemonitormodel.DeviceRunProcess, error)
}
