//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"google.golang.org/grpc"
	"nova-factory-server/app/business/metric/device/metricController"
	"nova-factory-server/app/business/metric/device/metricDao/metricDaoIMpl"
	"nova-factory-server/app/business/metric/device/metricService/metricServiceImpl"
	"nova-factory-server/app/datasource"
	"nova-factory-server/app/routes"
)

func wireApp() (*grpc.Server, func(), error) {
	panic(wire.Build(

		metricDaoIMpl.ProviderSet,
		metricServiceImpl.ProviderSet,
		metricController.ProviderSet,
		datasource.ProviderSet,
		routes.GrpcProviderSet))
}
