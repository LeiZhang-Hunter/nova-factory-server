package metricController

import (
	"context"
	"github.com/gin-gonic/gin"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"google.golang.org/grpc"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/business/metric/device/metricService"
	"nova-factory-server/app/utils/baizeContext"
)

type Metric struct {
	service             metricService.IMetricService
	deviceMetricService deviceMonitorService.DeviceMonitorService
	v1.UnimplementedDeviceReportServiceServer
}

func NewMetric(service metricService.IMetricService, deviceMetricService deviceMonitorService.DeviceMonitorService) *Metric {
	return &Metric{
		service:             service,
		deviceMetricService: deviceMetricService,
	}
}

func (m *Metric) PrivateGrpcRoutes(router *grpc.Server) {
	v1.RegisterDeviceReportServiceServer(router, m)
}

func (m *Metric) PrivateRoutes(router *gin.RouterGroup) {
}

func (m *Metric) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1")
	group.GET("/metric/predict", m.Predict)
}

func (m *Metric) ReportDeviceInfo(c context.Context, request *v1.ExportMetricsServiceRequest) (*v1.NodeRes, error) {
	err := m.service.Export(c, request)
	if err != nil {
		return &v1.NodeRes{
			Code: -1,
		}, err
	}
	return &v1.NodeRes{
		Code: 0,
	}, nil
}

// Predict 指标预测
// @Summary 导入告警数据
// @Description 导入告警数据
// @Tags 告警管理/告警数据管理
// @Param object body alertModels.AlertLogData true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /api/v1/predict/metric [post]
func (m *Metric) Predict(c *gin.Context) {
	req := new(metricModels.MetricQueryReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	data, err := m.deviceMetricService.Predict(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}
