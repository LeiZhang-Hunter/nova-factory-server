package deviceMonitorDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
)

type IDeviceDataReportDao interface {
	DevList(c *gin.Context) ([]deviceMonitorModel.SysIotDbDevMapData, error)
	GetDevList(c *gin.Context, dev []string) ([]deviceMonitorModel.SysIotDbDevMapData, error)
	Save(c *gin.Context, data *deviceMonitorModel.SysIotDbDevMap) error
	Remove(c *gin.Context, dev string) error
	List(c *gin.Context, req *deviceMonitorModel.DevListReq) (*deviceMonitorModel.DevListResp, error)
}
