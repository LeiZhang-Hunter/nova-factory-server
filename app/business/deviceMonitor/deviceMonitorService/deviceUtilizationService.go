package deviceMonitorService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
)

type DeviceUtilizationService interface {
	Stat(c *gin.Context, req *deviceMonitorModel.DeviceUtilizationReq) ([]*deviceMonitorModel.DeviceUtilizationData, error)
	Search(c *gin.Context, req *deviceMonitorModel.DeviceUtilizationReq) (*deviceMonitorModel.DeviceUtilizationPublicDataList, error)
	SearchV2(c *gin.Context, req *deviceMonitorModel.DeviceUtilizationReq) (*deviceMonitorModel.DeviceUtilizationPublicDataListV2, error)
}
