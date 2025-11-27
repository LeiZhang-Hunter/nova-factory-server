//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"nova-factory-server/app/business/ai/aiDataSetDao/aiDataSetDaoImpl"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorDao/deviceMonitorDaoImpl"
	"nova-factory-server/app/business/metric/device/metricDao/metricDaoIMpl"
	"nova-factory-server/app/datasource"
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
