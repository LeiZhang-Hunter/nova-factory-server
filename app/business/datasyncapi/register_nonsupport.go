//go:build !datasyncapi

// 数据同步 API 模块的降级注册（datasyncapi 编译标签未启用时）。
// 提供一个空的 DataSyncApi 和空的 ProviderSet，
// 避免在未启用该模块时因缺少依赖导致编译失败。
package datasyncapi

import (
	"github.com/google/wire"
)

// NewGinEngine 模块未启用时返回空的 DataSyncApi，不挂载任何路由
func NewGinEngine() *DataSyncApi {
	return &DataSyncApi{}
}

// GinProviderSet 模块未启用时的空 Gin 提供者集合
var GinProviderSet = wire.NewSet(NewGinEngine)

// ProviderSet 模块未启用时的空依赖注入集合
var ProviderSet = wire.NewSet(GinProviderSet)
