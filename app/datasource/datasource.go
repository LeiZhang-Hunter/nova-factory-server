package datasource

import (
	"github.com/google/wire"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/datasource/mysql"
	"nova-factory-server/app/datasource/objectFile"
)

var ProviderSet = wire.NewSet(mysql.NewData, mysql.NewDB, objectFile.NewConfig, cache.NewCache, cache.NewRedisSubscribe)
