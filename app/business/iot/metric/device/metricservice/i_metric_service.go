package metricservice

import (
	"context"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"

	"github.com/gin-gonic/gin"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
)

type IMetricService interface {
	Export(c context.Context, request *v1.ExportMetricsServiceRequest) error
	List(c *gin.Context, req *devicemonitormodel.DevDataReq) (*devicemonitormodel.DevDataResp, error)
	// ExportTimeData 导入时序数据
	ExportTimeData(c context.Context, request *v1.ExportTimeDataRequest) error

	ExportScheduleLog(c context.Context, request *v1.ExportTimeDataRequest) error
}
