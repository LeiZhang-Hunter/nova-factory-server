//go:build wireinject
// +build wireinject

package main

import (
	"nova-factory-server/app/business/ai/aiDataSetDao/aiDataSetDaoImpl"
	"nova-factory-server/app/business/iot/deviceMonitor/deviceMonitorDao/deviceMonitorDaoImpl"
	"nova-factory-server/app/business/iot/metric/device/metricDao/metricDaoIMpl"
	"nova-factory-server/app/datasource"

	"github.com/google/wire"
)

func wireApp() (*Runner, func(), error) {
	panic(wire.Build(
		deviceMonitorDaoImpl.ProviderSet,
		metricDaoIMpl.ProviderSet,
		aiDataSetDaoImpl.ProviderSet,
		datasource.ProviderSet,
		ProviderSet,
	))
}
