package metricdaoimpl

import (
	"context"
	"fmt"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/metric/device/metricdao"
	"nova-factory-server/app/business/iot/metric/device/metricmodels"
	"nova-factory-server/app/constant/datasource"
	"nova-factory-server/app/datasource/clickhouse"
	"nova-factory-server/app/datasource/iotdb"

	"github.com/gin-gonic/gin"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"github.com/spf13/viper"
)

// MetricDaoImpl 按配置的数据源委托指标读写实现。
type MetricDaoImpl struct {
	tableName string
	exporter  iDaoExport
}

func init() {

}

// NewMetricDaoImpl 根据 metric.datasource 创建指标 DAO。
func NewMetricDaoImpl() metricdao.IMetricDao {
	datasourceValue := viper.GetString("metric.datasource")
	var exporter iDaoExport
	switch datasourceValue {
	case datasource.IOTDB:
		{
			exporter = newIotDbExport(iotdb.GetIotDb())
			break
		}
	case datasource.CLICKHOUSE:
		{
			exporter = newIClickHouseExport(clickhouse.GetClickHouse())
			break
		}
	case datasource.PROMETHEUS_TSDB:
		{
			exporter = newPrometheusTSDBExport(iotdb.GetTSDBStorage())
		}
	default:
		panic(fmt.Sprintf("datasource: %s is not exist", datasourceValue))
	}
	return &MetricDaoImpl{
		exporter: exporter,
	}
}

// Export 写入设备指标数据。
func (m *MetricDaoImpl) Export(ctx context.Context, data []*metricmodels.NovaMetricsDevice) error {
	return m.exporter.Export(ctx, data)
}

// Metric 查询单个设备测点的聚合指标。
func (m *MetricDaoImpl) Metric(c *gin.Context, req *metricmodels.MetricQueryReq) (*metricmodels.MetricQueryData, error) {
	return m.exporter.Metric(c, req)
}

// CounterByTimeRange 按时间范围统计设备写入数量。
func (m *MetricDaoImpl) CounterByTimeRange(startTime int64, endTime int64, interval string) (*metricmodels.MetricQueryData, error) {
	exporter, ok := m.exporter.(interface {
		CounterByTimeRange(startTime int64, endTime int64, interval string) (*metricmodels.MetricQueryData, error)
	})
	if !ok {
		return metricmodels.NewMetricQueryData(), nil
	}
	return exporter.CounterByTimeRange(startTime, endTime, interval)
}

// CounterByDevice 按设备维度统计写入数量排行。
func (m *MetricDaoImpl) CounterByDevice(c *gin.Context, startTime int64, endTime int64, limit int) (*devicemonitormodel.TypeDeviceCounterRank, error) {
	exporter, ok := m.exporter.(interface {
		CounterByDevice(c *gin.Context, startTime int64, endTime int64, limit int) (*devicemonitormodel.TypeDeviceCounterRank, error)
	})
	if !ok {
		return &devicemonitormodel.TypeDeviceCounterRank{Rows: make([]*devicemonitormodel.TypeDeviceCounterRankValue, 0)}, nil
	}
	return exporter.CounterByDevice(c, startTime, endTime, limit)
}

// StatDeviceStatus 按状态统计设备运行时长。
func (m *MetricDaoImpl) StatDeviceStatus(c *gin.Context, startTime string, endTime string, status int) (*devicemonitormodel.DeviceStatusList, error) {
	exporter, ok := m.exporter.(interface {
		StatDeviceStatus(c *gin.Context, startTime string, endTime string, status int) (*devicemonitormodel.DeviceStatusList, error)
	})
	if !ok {
		return devicemonitormodel.NewDeviceStatusList(), nil
	}
	return exporter.StatDeviceStatus(c, startTime, endTime, status)
}

// StatDeviceProcess 按时间分组统计设备运行过程。
func (m *MetricDaoImpl) StatDeviceProcess(c *gin.Context, startTime string, endTime string, interval string, status int) (*devicemonitormodel.DeviceProcessList, error) {
	exporter, ok := m.exporter.(interface {
		StatDeviceProcess(c *gin.Context, startTime string, endTime string, interval string, status int) (*devicemonitormodel.DeviceProcessList, error)
	})
	if !ok {
		return &devicemonitormodel.DeviceProcessList{List: make(map[string][]devicemonitormodel.DeviceStatus)}, nil
	}
	return exporter.StatDeviceProcess(c, startTime, endTime, interval, status)
}

// StatDeviceRunStatus 查询设备在时间范围内的最后运行状态。
func (m *MetricDaoImpl) StatDeviceRunStatus(c *gin.Context, startTime string, endTime string) ([]devicemonitormodel.DeviceRunStat, error) {
	exporter, ok := m.exporter.(interface {
		StatDeviceRunStatus(c *gin.Context, startTime string, endTime string) ([]devicemonitormodel.DeviceRunStat, error)
	})
	if !ok {
		return make([]devicemonitormodel.DeviceRunStat, 0), nil
	}
	return exporter.StatDeviceRunStatus(c, startTime, endTime)
}

// StatDeviceStatusByDeviceId 按设备和状态统计运行时长。
func (m *MetricDaoImpl) StatDeviceStatusByDeviceId(c *gin.Context, startTime string, endTime string, deviceId int64, status int) (*devicemonitormodel.DeviceStatusList, error) {
	exporter, ok := m.exporter.(interface {
		StatDeviceStatusByDeviceId(c *gin.Context, startTime string, endTime string, deviceId int64, status int) (*devicemonitormodel.DeviceStatusList, error)
	})
	if !ok {
		return devicemonitormodel.NewDeviceStatusList(), nil
	}
	return exporter.StatDeviceStatusByDeviceId(c, startTime, endTime, deviceId, status)
}

// StatDeviceProcessByDeviceId 按设备统计运行过程。
func (m *MetricDaoImpl) StatDeviceProcessByDeviceId(c *gin.Context, startTime string, endTime string, deviceId int64, interval string, status int) (*devicemonitormodel.DeviceProcessList, error) {
	exporter, ok := m.exporter.(interface {
		StatDeviceProcessByDeviceId(c *gin.Context, startTime string, endTime string, deviceId int64, interval string, status int) (*devicemonitormodel.DeviceProcessList, error)
	})
	if !ok {
		return &devicemonitormodel.DeviceProcessList{List: make(map[string][]devicemonitormodel.DeviceStatus)}, nil
	}
	return exporter.StatDeviceProcessByDeviceId(c, startTime, endTime, deviceId, interval, status)
}

// Predict 查询设备指标预测数据。
func (m *MetricDaoImpl) Predict(c *gin.Context, deviceId int64, device *devicemodels.SysModbusDeviceConfigData, req *metricmodels.MetricQueryReq) (*metricmodels.MetricQueryData, error) {
	return m.exporter.Predict(c, deviceId, device, req)
}

// InstallDevice 安装设备指标存储结构。
func (m *MetricDaoImpl) InstallDevice(c *gin.Context, deviceId int64, device *devicemodels.SysModbusDeviceConfigData) error {
	return m.exporter.InstallDevice(c, deviceId, device)
}

// UnInStallDevice 卸载设备指标存储结构。
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

// List 查询设备原始时序数据列表。
func (m *MetricDaoImpl) List(c *gin.Context, req *devicemonitormodel.DevDataReq) (*devicemonitormodel.DevDataResp, error) {
	return m.exporter.List(c, req)
}

// Count 统计设备原始时序数据数量。
func (m *MetricDaoImpl) Count(c *gin.Context, req *devicemonitormodel.DevDataReq) (uint64, error) {
	return m.exporter.Count(c, req)
}

// Query 查询 dashboard 指标数据。
func (m *MetricDaoImpl) Query(c *gin.Context, req *metricmodels.MetricDataQueryReq) (*metricmodels.MetricQueryData, error) {
	return m.exporter.Query(c, req)
}

// ExportTimeData 导出时序数据
func (m *MetricDaoImpl) ExportTimeData(ctx context.Context, data map[string][]*v1.ResourceTimeMetrics) error {
	return m.exporter.ExportTimeData(ctx, data)
}
