package devicemonitorservice

import (
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/metric/device/metricmodels"

	"github.com/gin-gonic/gin"
)

type ControlLogService interface {
	// List 控制日志列表
	List(c *gin.Context, req *devicemonitormodel.ControlLogListReq) (*metricmodels.NovaControlLogList, error)
}
