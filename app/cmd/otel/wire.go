//go:build wireinject
// +build wireinject

package main

import (
	"nova-factory-server/app/business/iot/asset/camera/cameradao/cameraDaoImpl"
	"nova-factory-server/app/business/iot/metric/device/metriccontroller"
	"nova-factory-server/app/business/iot/metric/device/metricdao/metricdaoimpl"
	"nova-factory-server/app/business/iot/metric/device/metricservice/metricserviceimpl"
	"nova-factory-server/app/datasource"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
	"google.golang.org/grpc"
)

func wireApp() (*grpc.Server, func(), error) {
	panic(wire.Build(

		metricdaoimpl.ProviderSet,
		metricserviceimpl.ProviderSet,
		metriccontroller.ProviderSet,
		datasource.ProviderSet,
		cameraDaoImpl.ProviderSet,
		routes.GrpcProviderSet))
}
