package clickhouse

import (
	"context"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"go.uber.org/zap"
	metricmodels "nova-factory-server/app/business/iot/metric/device/metricmodels/entity"
	"nova-factory-server/app/datasource/clickhouse"
)

type adder struct {
	tableName  string
	clickhouse *clickhouse.ClickHouse
}

func newAdder(tableName string, clickhouse *clickhouse.ClickHouse) *adder {
	return &adder{
		tableName:  tableName,
		clickhouse: clickhouse,
	}
}

func (a *adder) Export(ctx context.Context, data []*metricmodels.NovaMetricsDevice) error {
	if len(data) == 0 {
		return nil
	}
	ret := a.clickhouse.DB().Table(a.tableName).Debug().Create(data)
	if ret.Error != nil {
		zap.L().Error("create device metric data error:", zap.Error(ret.Error))
		return ret.Error
	}
	return ret.Error
}

func (a *adder) ExportTimeData(ctx context.Context, data map[string][]*v1.ResourceTimeMetrics) error {
	return nil
}
