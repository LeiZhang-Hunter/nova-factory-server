package deviceMonitorDao

import (
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"

	"github.com/gin-gonic/gin"
)

type DeviceUtilizationDao interface {
	Stat(c *gin.Context, req *deviceMonitorModel.DeviceUtilizationReq) ([]*deviceMonitorModel.DeviceUtilizationData, error)
	Search(c *gin.Context, req *deviceMonitorModel.DeviceUtilizationReq) (*deviceMonitorModel.DeviceUtilizationPublicDataList, error)
	SearchV2(c *gin.Context, req *deviceMonitorModel.DeviceUtilizationReq) (*deviceMonitorModel.DeviceUtilizationPublicDataListV2, error)
	GetDeviceUtilization(c *gin.Context, req *deviceMonitorModel.DeviceUtilizationReq) (*deviceMonitorModel.DeviceRunProcess, error)
}
