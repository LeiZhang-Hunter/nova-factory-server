package data

import (
	"nova-factory-server/app/business/data/agent"
	"nova-factory-server/app/business/data/controller"
	"nova-factory-server/app/business/data/dao/impl"
	serviceimpl "nova-factory-server/app/business/data/service/impl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	impl.ProviderSet,
	serviceimpl.ProviderSet,
	agent.ProviderSet,
	controller.ProviderSet,

	GinProviderSet,
)
