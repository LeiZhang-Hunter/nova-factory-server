package metricDaoIMpl

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"github.com/spf13/viper"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/metric/device/metricDao"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/constant/datasource"
	"nova-factory-server/app/datasource/clickhouse"
	"nova-factory-server/app/datasource/iotdb"
)

type MetricDaoImpl struct {
	tableName string
	exporter  iDaoExport
}

func init() {

}

func NewMetricDaoImpl(clickhouse *clickhouse.ClickHouse,
	iotDb *iotdb.IotDb) metricDao.IMetricDao {
	datasourceValue := viper.GetString("metric.datasource")
	var exporter iDaoExport
	switch datasourceValue {
	case datasource.IOTDB:
		{
			exporter = newIotDbExport(iotDb)
			break
		}
	case datasource.CLICKHOUSE:
		{
			exporter = newIClickHouseExport(clickhouse)
			break
		}
	default:
		panic(fmt.Sprintf("datasource: %s is not exist", datasourceValue))
	}
	return &MetricDaoImpl{
		exporter: exporter,
	}
}

func (m *MetricDaoImpl) Export(ctx context.Context, data []*metricModels.NovaMetricsDevice) error {
	return m.exporter.Export(ctx, data)
}

func (m *MetricDaoImpl) Metric(c *gin.Context, req *metricModels.MetricQueryReq) (*metricModels.MetricQueryData, error) {
	return m.exporter.Metric(c, req)
}

func (m *MetricDaoImpl) Predict(c *gin.Context, deviceId int64, device *deviceModels.SysModbusDeviceConfigData, req *metricModels.MetricQueryReq) (*metricModels.MetricQueryData, error) {
	return m.exporter.Predict(c, deviceId, device, req)
}

func (m *MetricDaoImpl) InstallDevice(c *gin.Context, deviceId int64, device *deviceModels.SysModbusDeviceConfigData) error {
	return m.exporter.InstallDevice(c, deviceId, device)
}

func (m *MetricDaoImpl) UnInStallDevice(c *gin.Context, deviceId int64, templateId int64, dataId int64) error {
	return m.exporter.UnInStallDevice(c, deviceId, templateId, dataId)
}

// InstallRunStatusDevice 运行状态设备模板
func (m *MetricDaoImpl) InstallRunStatusDevice(c *gin.Context, deviceId int64) error {
	return m.exporter.InstallRunStatusDevice(c, deviceId)
}

// UnInStallRunStatusDevice 卸载设备运行状态模板
func (m *MetricDaoImpl) UnInStallRunStatusDevice(c *gin.Context, deviceId int64) error {
	return m.exporter.UnInStallRunStatusDevice(c, deviceId)
}

func (m *MetricDaoImpl) List(c *gin.Context, req *deviceMonitorModel.DevDataReq) (*deviceMonitorModel.DevDataResp, error) {
	return m.exporter.List(c, req)
}

func (m *MetricDaoImpl) Count(c *gin.Context, req *deviceMonitorModel.DevDataReq) (uint64, error) {
	return m.exporter.Count(c, req)
}

func (m *MetricDaoImpl) Query(c *gin.Context, req *metricModels.MetricDataQueryReq) (*metricModels.MetricQueryData, error) {
	return m.exporter.Query(c, req)
}

// ExportTimeData 导出时序数据
func (m *MetricDaoImpl) ExportTimeData(ctx context.Context, data map[string][]*v1.ResourceTimeMetrics) error {
	return m.exporter.ExportTimeData(ctx, data)
}
