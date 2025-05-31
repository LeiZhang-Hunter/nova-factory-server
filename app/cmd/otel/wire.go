//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"google.golang.org/grpc"
	"nova-factory-server/app/business/metric/device/metricController"
	"nova-factory-server/app/routes"
)

func wireApp() (*grpc.Server, func(), error) {
	panic(wire.Build(metricController.ProviderSet, routes.GrpcProviderSet))
}
