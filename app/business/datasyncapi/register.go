//go:build datasyncapi
// +build datasyncapi

// 数据同步 API 模块的依赖注入注册（datasyncapi 编译标签启用时）。
// 通过 wire 汇总管家婆控制器和服务的 ProviderSet，
// 以及 Gin 引擎的 ProviderSet，向全局注入器提供完整依赖链。
package datasyncapi

import (
	qqdcontroller "nova-factory-server/app/business/datasyncapi/gjpqqd/controller"
	qqdservice "nova-factory-server/app/business/datasyncapi/gjpqqd/service/impl"

	"github.com/google/wire"
)

// ProviderSet 数据同步 API 的 wire 依赖注入集合，
// 按需汇总子模块（控制器、服务实现、Gin 引擎）的 ProviderSet
var ProviderSet = wire.NewSet(
	qqdcontroller.ProviderSet,
	qqdservice.ProviderSet,
	GinProviderSet,
)
