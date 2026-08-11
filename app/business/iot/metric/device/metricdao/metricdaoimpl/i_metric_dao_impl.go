package metricdaoimpl

import (
	"context"
	"fmt"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/metric/device/metricdao"
	clickhouseimpl "nova-factory-server/app/business/iot/metric/device/metricdao/metricdaoimpl/clickhouse"
	iotdbimpl "nova-factory-server/app/business/iot/metric/device/metricdao/metricdaoimpl/iotdb"
	prometheusimpl "nova-factory-server/app/business/iot/metric/device/metricdao/metricdaoimpl/prometheus_tsdb"
	"nova-factory-server/app/business/iot/metric/device/metricmodels/entity"
	"nova-factory-server/app/constant/datasource"
	"nova-factory-server/app/datasource/clickhouse"
	"nova-factory-server/app/datasource/iotdb"

	"github.com/gin-gonic/gin"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"github.com/spf13/viper"
)

// MetricDaoImpl delegates metric read/write operations to the configured storage backend.
type MetricDaoImpl struct {
	storage metricdao.IMetricStorageDao
}

func NewMetricDaoImpl() metricdao.IMetricDao {
	datasourceValue := viper.GetString("metric.datasource")
	var storage metricdao.IMetricStorageDao
	switch datasourceValue {
	case datasource.IOTDB:
		storage = iotdbimpl.NewIotDbExport(iotdb.GetIotDb())
	case datasource.CLICKHOUSE:
		storage = clickhouseimpl.NewIClickHouseExport(clickhouse.GetClickHouse())
	case datasource.PROMETHEUS_TSDB:
		storage = prometheusimpl.NewPrometheusTSDBExport(iotdb.GetTSDBStorage())
	default:
		panic(fmt.Sprintf("datasource: %s is not exist", datasourceValue))
	}
	return &MetricDaoImpl{storage: storage}
}

func (m *MetricDaoImpl) Export(ctx context.Context, data []*entity.NovaMetricsDevice) error {
	if m.storage.Adder() == nil {
		return nil
	}
	return m.storage.Adder().Export(ctx, data)
}

func (m *MetricDaoImpl) ExportTimeData(ctx context.Context, data map[string][]*v1.ResourceTimeMetrics) error {
	if m.storage.Adder() == nil {
		return nil
	}
	return m.storage.Adder().ExportTimeData(ctx, data)
}

func (m *MetricDaoImpl) Metric(c *gin.Context, req *entity.MetricQueryReq) (*entity.MetricQueryData, error) {
	if m.storage.Questioner() == nil {
		return nil, nil
	}
	return m.storage.Questioner().Metric(c, req)
}

func (m *MetricDaoImpl) Predict(c *gin.Context, deviceId int64, device *devicemodels.SysModbusDeviceConfigData, req *entity.MetricQueryReq) (*entity.MetricQueryData, error) {
	if m.storage.Questioner() == nil {
		return nil, nil
	}
	return m.storage.Questioner().Predict(c, deviceId, device, req)
}

func (m *MetricDaoImpl) List(c *gin.Context, req *devicemonitormodel.DevDataReq) (*devicemonitormodel.DevDataResp, error) {
	if m.storage.Questioner() == nil {
		return nil, nil
	}
	return m.storage.Questioner().List(c, req)
}

func (m *MetricDaoImpl) Count(c *gin.Context, req *devicemonitormodel.DevDataReq) (uint64, error) {
	if m.storage.Questioner() == nil {
		return 0, nil
	}
	return m.storage.Questioner().Count(c, req)
}

func (m *MetricDaoImpl) Query(c *gin.Context, req *entity.MetricDataQueryReq) (*entity.MetricQueryData, error) {
	if m.storage.Questioner() == nil {
		return nil, nil
	}
	return m.storage.Questioner().Query(c, req)
}

func (m *MetricDaoImpl) InstallDevice(c *gin.Context, deviceId int64, templateId int64) error {
	return m.storage.InstallDevice(c, deviceId, templateId)
}

func (m *MetricDaoImpl) UnInStallDevice(c *gin.Context, deviceId int64, templateId int64) error {
	return m.storage.UnInStallDevice(c, deviceId, templateId)
}

func (m *MetricDaoImpl) InstallRunStatusDevice(c *gin.Context, deviceId int64) error {
	return m.storage.InstallRunStatusDevice(c, deviceId)
}

func (m *MetricDaoImpl) UnInStallRunStatusDevice(c *gin.Context, deviceId int64) error {
	return m.storage.UnInStallRunStatusDevice(c, deviceId)
}

func (m *MetricDaoImpl) Template() metricdao.IIotStorageTemplateDao {
	return m.storage.Template()
}

func (m *MetricDaoImpl) Adder() metricdao.IMetricAdderDao {
	if adder := m.storage.Adder(); adder != nil {
		return adder
	}
	return m.storage.Adder()
}

func (m *MetricDaoImpl) CounterByTimeRange(startTime int64, endTime int64, interval string) (*entity.MetricQueryData, error) {
	storage, ok := m.storage.(interface {
		CounterByTimeRange(startTime int64, endTime int64, interval string) (*entity.MetricQueryData, error)
	})
	if !ok {
		return entity.NewMetricQueryData(), nil
	}
	return storage.CounterByTimeRange(startTime, endTime, interval)
}

func (m *MetricDaoImpl) CounterByDevice(c *gin.Context, startTime int64, endTime int64, limit int) (*devicemonitormodel.TypeDeviceCounterRank, error) {
	storage, ok := m.storage.(interface {
		CounterByDevice(c *gin.Context, startTime int64, endTime int64, limit int) (*devicemonitormodel.TypeDeviceCounterRank, error)
	})
	if !ok {
		return &devicemonitormodel.TypeDeviceCounterRank{Rows: make([]*devicemonitormodel.TypeDeviceCounterRankValue, 0)}, nil
	}
	return storage.CounterByDevice(c, startTime, endTime, limit)
}

func (m *MetricDaoImpl) StatDeviceStatus(c *gin.Context, startTime string, endTime string, status int) (*devicemonitormodel.DeviceStatusList, error) {
	storage, ok := m.storage.(interface {
		StatDeviceStatus(c *gin.Context, startTime string, endTime string, status int) (*devicemonitormodel.DeviceStatusList, error)
	})
	if !ok {
		return devicemonitormodel.NewDeviceStatusList(), nil
	}
	return storage.StatDeviceStatus(c, startTime, endTime, status)
}

func (m *MetricDaoImpl) StatDeviceProcess(c *gin.Context, startTime string, endTime string, interval string, status int) (*devicemonitormodel.DeviceProcessList, error) {
	storage, ok := m.storage.(interface {
		StatDeviceProcess(c *gin.Context, startTime string, endTime string, interval string, status int) (*devicemonitormodel.DeviceProcessList, error)
	})
	if !ok {
		return &devicemonitormodel.DeviceProcessList{List: make(map[string][]devicemonitormodel.DeviceStatus)}, nil
	}
	return storage.StatDeviceProcess(c, startTime, endTime, interval, status)
}

func (m *MetricDaoImpl) StatDeviceRunStatus(c *gin.Context, startTime string, endTime string) ([]devicemonitormodel.DeviceRunStat, error) {
	storage, ok := m.storage.(interface {
		StatDeviceRunStatus(c *gin.Context, startTime string, endTime string) ([]devicemonitormodel.DeviceRunStat, error)
	})
	if !ok {
		return make([]devicemonitormodel.DeviceRunStat, 0), nil
	}
	return storage.StatDeviceRunStatus(c, startTime, endTime)
}

func (m *MetricDaoImpl) StatDeviceStatusByDeviceId(c *gin.Context, startTime string, endTime string, deviceId int64, status int) (*devicemonitormodel.DeviceStatusList, error) {
	storage, ok := m.storage.(interface {
		StatDeviceStatusByDeviceId(c *gin.Context, startTime string, endTime string, deviceId int64, status int) (*devicemonitormodel.DeviceStatusList, error)
	})
	if !ok {
		return devicemonitormodel.NewDeviceStatusList(), nil
	}
	return storage.StatDeviceStatusByDeviceId(c, startTime, endTime, deviceId, status)
}

func (m *MetricDaoImpl) StatDeviceProcessByDeviceId(c *gin.Context, startTime string, endTime string, deviceId int64, interval string, status int) (*devicemonitormodel.DeviceProcessList, error) {
	storage, ok := m.storage.(interface {
		StatDeviceProcessByDeviceId(c *gin.Context, startTime string, endTime string, deviceId int64, interval string, status int) (*devicemonitormodel.DeviceProcessList, error)
	})
	if !ok {
		return &devicemonitormodel.DeviceProcessList{List: make(map[string][]devicemonitormodel.DeviceStatus)}, nil
	}
	return storage.StatDeviceProcessByDeviceId(c, startTime, endTime, deviceId, interval, status)
}
