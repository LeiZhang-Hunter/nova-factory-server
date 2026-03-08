//go:build wireinject
// +build wireinject

package main

import (
	"nova-factory-server/app/business/iot/metric/device/metricController"
	"nova-factory-server/app/business/iot/metric/device/metricDao/metricDaoIMpl"
	"nova-factory-server/app/business/iot/metric/device/metricService/metricServiceImpl"
	"nova-factory-server/app/datasource"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
	"google.golang.org/grpc"
)

func wireApp() (*grpc.Server, func(), error) {
	panic(wire.Build(

		metricDaoIMpl.ProviderSet,
		metricServiceImpl.ProviderSet,
		metricController.ProviderSet,
		datasource.ProviderSet,
		routes.GrpcProviderSet))
}
