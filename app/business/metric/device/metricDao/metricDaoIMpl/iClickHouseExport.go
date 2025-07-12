package metricDaoIMpl

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/datasource/clickhouse"
	"nova-factory-server/app/utils/time"
)

type iClickHouseExport struct {
	tableName  string
	clickhouse *clickhouse.ClickHouse
}

func newIClickHouseExport(clickhouse *clickhouse.ClickHouse) iDaoExport {
	return &iClickHouseExport{
		clickhouse: clickhouse,
		tableName:  "nova_metrics_device",
	}
}

func (m *iClickHouseExport) Export(ctx context.Context, data []*metricModels.NovaMetricsDevice) error {
	if len(data) == 0 {
		return nil
	}
	ret := m.clickhouse.DB().Table(m.tableName).Debug().Create(data)
	if ret.Error != nil {
		zap.L().Error("create device metric data error:", zap.Error(ret.Error))
		return ret.Error
	}
	return ret.Error
}

func (m *iClickHouseExport) Metric(c *gin.Context, req *metricModels.MetricQueryReq) (*metricModels.MetricQueryData, error) {
	if req == nil {
		return nil, nil
	}

	model := m.clickhouse.DB().Table(m.tableName)

	var startTime string
	if req.Start > 0 {
		startTime = time.GetStartTime(req.Start, 200)
	}
	endTime := time.GetEndTimeUseNow(req.End, true)

	if startTime != "" && endTime != "" {

		model = model.Where("time_unix >= ?", startTime)
	}
	if endTime != "" {
		model = model.Where("time_unix <= ?", endTime)
	}
	if req.Step <= 0 {
		req.Step = 1
	}

	var list []*metricModels.NovaMetricsDevice = make([]*metricModels.NovaMetricsDevice, 0)
	model = model.
		Select(fmt.Sprintf(
			"`toStartOfInterval`(`time_unix`, INTERVAL %d minute) AS `time_unix`, %s as value",
			req.Step, "AVG(value)"))

	ret := model.Where("device_id = ?", req.DeviceId).
		Where("template_id = ?", req.TemplateId).
		Where("data_id = ?", req.DataId).
		Group("time_unix").Order("time_unix asc").Limit(500).Find(&list)
	if ret.Error != nil {
		return nil, ret.Error
	}

	if len(list) == 0 {
		return &metricModels.MetricQueryData{
			Labels: make(map[string]string),
			Values: make([]metricModels.MetricQueryValue, len(list)),
		}, nil
	}

	var data metricModels.MetricQueryData
	data.Labels = make(map[string]string)
	data.Values = make([]metricModels.MetricQueryValue, len(list))
	for i := 0; i < len(list); i++ {
		data.Values[i].Time = list[i].TimeUnix.UnixMilli()
		data.Values[i].Value = fmt.Sprintf("%f", list[i].Value)
	}
	return &data, nil
}

func (m *iClickHouseExport) InstallDevice(c *gin.Context, deviceId int64, device *deviceModels.SysModbusDeviceConfigData) error {
	return nil
}
func (m *iClickHouseExport) UnInStallDevice(c *gin.Context, deviceId int64, templateId int64, dataId int64) error {
	return nil
}

func (i *iClickHouseExport) Predict(c *gin.Context, deviceId int64, device *deviceModels.SysModbusDeviceConfigData, req *metricModels.MetricQueryReq) (*metricModels.MetricQueryData, error) {
	return &metricModels.MetricQueryData{
		Labels: make(map[string]string),
		Values: make([]metricModels.MetricQueryValue, 0),
	}, nil
}

func (i *iClickHouseExport) List(c *gin.Context, req *deviceMonitorModel.DevDataReq) (*deviceMonitorModel.DevDataResp, error) {
	return nil, nil
}
