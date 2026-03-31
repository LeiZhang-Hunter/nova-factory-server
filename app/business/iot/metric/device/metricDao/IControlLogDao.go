package metricDao

import (
	"context"
	"nova-factory-server/app/business/iot/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/iot/metric/device/metricModels"

	"github.com/gin-gonic/gin"
)

type IControlLogDao interface {
	// Export 导出数据
	Export(ctx context.Context, data []*metricModels.NovaControlLog) error
	// List 控制日志列表
	List(c *gin.Context, req *deviceMonitorModel.ControlLogListReq) (*metricModels.NovaControlLogList, error)
}
