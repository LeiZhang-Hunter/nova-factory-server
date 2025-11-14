package deviceMonitorService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/metric/device/metricModels"
)

type ControlLogService interface {
	// List 控制日志列表
	List(c *gin.Context, req *deviceMonitorModel.ControlLogListReq) (*metricModels.NovaControlLogList, error)
}
