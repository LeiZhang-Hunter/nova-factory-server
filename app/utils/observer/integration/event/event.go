package event

import (
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/config"
)

type Base interface {
	Metadata() map[string]any
	Ptr() any
}

type Event interface {
	Config() config.Config
	Action() EventType
	Cache() cache.Cache
}
