package metricDaoIMpl

import (
	"context"
	"go.uber.org/zap"
	"nova-factory-server/app/business/metric/device/metricDao"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/datasource/clickhouse"
)

type MetricDaoImpl struct {
	clickhouse *clickhouse.ClickHouse
	tableName  string
}

func NewMetricDaoImpl(clickhouse *clickhouse.ClickHouse) metricDao.IMetricDao {
	return &MetricDaoImpl{
		clickhouse: clickhouse,
		tableName:  "nova_metrics_device",
	}
}

func (m *MetricDaoImpl) Export(ctx context.Context, data []*metricModels.NovaMetricsDevice) error {
	if len(data) == 0 {
		return nil
	}
	ret := m.clickhouse.DB().Table(m.tableName).Create(data)
	if ret.Error != nil {
		zap.L().Error("create device metric data error:", zap.Error(ret.Error))
		return ret.Error
	}
	return ret.Error
}
