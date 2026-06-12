//go:build datasyncapi
// +build datasyncapi

// 数据同步 API 的路由注册与 Gin 引擎组装。
// 仅在编译标签 datasyncapi 启用时编译，
// 将管家婆全渠道（gjpqqd）控制器的公开路由挂载到 /api/v1/erp-api/qqd 下。
package datasyncapi

import (
	qqdcontroller "nova-factory-server/app/business/datasyncapi/gjpqqd/controller"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

// GinProviderSet 提供 Gin 引擎的 wire 依赖集合
var GinProviderSet = wire.NewSet(NewGinEngine)

// NewGinEngine 创建数据同步 API 的 Gin 路由组
// 将 /api/v1/erp-api/qqd 前缀下的路由挂载到管家婆控制器
func NewGinEngine(app *routes.App, qqd *qqdcontroller.Controller) *DataSyncApi {
	group := app.Engine.Group("/api/v1/erp-api")
	qqd.PublicRoutes(group)
	return &DataSyncApi{}
}
