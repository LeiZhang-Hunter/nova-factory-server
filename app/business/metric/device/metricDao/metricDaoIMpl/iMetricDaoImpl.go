package metricDaoIMpl

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"nova-factory-server/app/business/asset/device/deviceModels"
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

func (m *MetricDaoImpl) InstallDevice(c *gin.Context, deviceId int64, device *deviceModels.SysModbusDeviceConfigData) error {
	return m.exporter.InstallDevice(c, deviceId, device)
}

func (m *MetricDaoImpl) UnInStallDevice(c *gin.Context, deviceId int64, templateId int64, dataId int64) error {
	return m.exporter.UnInStallDevice(c, deviceId, templateId, dataId)
}
