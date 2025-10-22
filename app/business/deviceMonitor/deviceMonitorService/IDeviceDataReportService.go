package deviceMonitorService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
)

type IDeviceDataReportService interface {
	DevList(c *gin.Context) ([]deviceMonitorModel.SysIotDbDevMapData, error)
	GetDevList(c *gin.Context, req *deviceMonitorModel.DevListReq) (*deviceMonitorModel.DevListResp, error)
	// SetDevMap 设置测点
	SetDevMap(c *gin.Context, info *deviceMonitorModel.SetDevMapInfo) error
	// RemoveDevMap 删除测点
	RemoveDevMap(c *gin.Context, dev string) error
}
