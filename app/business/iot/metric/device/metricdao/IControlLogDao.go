package metricdao

import (
	"context"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/metric/device/metricmodels"

	"github.com/gin-gonic/gin"
)

type IControlLogDao interface {
	// Export 导出数据
	Export(ctx context.Context, data []*metricmodels.NovaControlLog) error
	// List 控制日志列表
	List(c *gin.Context, req *devicemonitormodel.ControlLogListReq) (*metricmodels.NovaControlLogList, error)
}
