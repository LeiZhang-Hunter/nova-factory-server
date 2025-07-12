package metricService

import (
	"context"
	"github.com/gin-gonic/gin"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
)

type IMetricService interface {
	Export(c context.Context, request *v1.ExportMetricsServiceRequest) error
	List(c *gin.Context, req *deviceMonitorModel.DevDataReq) (*deviceMonitorModel.DevDataResp, error)
}
