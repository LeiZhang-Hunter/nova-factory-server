//go:build wireinject
// +build wireinject

package main

import (
	"nova-factory-server/app/business/admin"
	"nova-factory-server/app/business/ai"
	"nova-factory-server/app/business/datasyncapi"
	"nova-factory-server/app/business/erp"
	"nova-factory-server/app/business/iot"
	"nova-factory-server/app/business/shop"
	"nova-factory-server/app/business/wms"
	"nova-factory-server/app/datasource"
	"nova-factory-server/app/routes"
	"nova-factory-server/app/utils/observer/integration/observer"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

func finalEngine(app *routes.App, _ *admin.Admin, _ *iot.Iot, _ *ai.AI, _ *shop.Shop, _ *erp.Erp, _ *datasyncapi.DataSyncApi, _ *wms.Wms, _ *observer.Notifier) *gin.Engine {
	type McpConfig struct {
		Path           string `mapstructure:"path"`
		OperationsPath string `mapstructure:"operationsPath"`
	}
	// 把读取到的配置信息反序列化到 Conf 变量中
	var mcpConfig McpConfig
	if err := viper.UnmarshalKey("mcp", &mcpConfig); err != nil {
		panic(err)
	}

	// 4. Mount the MCP server endpoint
	app.McpServer.Mount() // MCP clients will connect here
	app.GrpcServer.Start()
	return app.Engine
}

func wireApp() (*gin.Engine, func(), error) {
	panic(wire.Build(
		routes.ProviderSet,
		iot.ProviderSet,
		ai.ProviderSet,
		shop.ProviderSet,
		erp.ProviderSet,
		datasyncapi.ProviderSet,
		wms.ProviderSet,
		admin.ProviderSet,
		datasource.ProviderSet,
		finalEngine,
	))
}
