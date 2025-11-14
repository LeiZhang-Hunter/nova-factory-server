package metricDaoIMpl

import (
	"context"
	"go.uber.org/zap"
	"nova-factory-server/app/business/metric/device/metricDao"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/datasource/clickhouse"
)

type IControlLogDaoImpl struct {
	tableName  string
	clickhouse *clickhouse.ClickHouse
}

func NewIControlLogDaoImpl(clickhouse *clickhouse.ClickHouse) metricDao.IControlLogDao {
	return &IControlLogDaoImpl{
		clickhouse: clickhouse,
		tableName:  "nova_control_log",
	}
}

func (i *IControlLogDaoImpl) Export(ctx context.Context, data []*metricModels.NovaControlLog) error {
	if len(data) == 0 {
		return nil
	}
	ret := i.clickhouse.DB().Table(i.tableName).Create(data)
	if ret.Error != nil {
		zap.L().Error("create device metric data error:", zap.Error(ret.Error))
		return ret.Error
	}
	return ret.Error
}
