package deviceMonitorService

import (
	"nova-factory-server/app/business/iot/deviceMonitor/deviceMonitorModel"

	"github.com/gin-gonic/gin"
)

type IDeviceDataReportService interface {
	DevList(c *gin.Context) ([]deviceMonitorModel.SysIotDbDevMapData, error)
	GetDevList(c *gin.Context, req *deviceMonitorModel.DevListReq) (*deviceMonitorModel.DevListResp, error)
	// SetDevMap 设置测点
	SetDevMap(c *gin.Context, info *deviceMonitorModel.SetDevMapInfo) error
	// RemoveDevMap 删除测点
	RemoveDevMap(c *gin.Context, dev string) error
}
