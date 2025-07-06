package datasource

import (
	"github.com/google/wire"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/datasource/clickhouse"
	"nova-factory-server/app/datasource/iotdb"
	"nova-factory-server/app/datasource/mysql"
	"nova-factory-server/app/datasource/objectFile"
)

var ProviderSet = wire.NewSet(mysql.NewData, mysql.NewDB, clickhouse.NewClickHouse, iotdb.NewIotDb, objectFile.NewConfig, cache.NewCache, cache.NewRedisSubscribe)
