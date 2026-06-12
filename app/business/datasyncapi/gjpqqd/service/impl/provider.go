// 服务实现层的 wire 依赖注入注册。
// 通过 ProviderSet 导出 IQQDServiceImpl 的构造函数，供上层注入器使用。
package impl

import "github.com/google/wire"

// ProviderSet 服务实现的 wire 依赖集合，提供 NewIQQDServiceImpl 构造函数
var ProviderSet = wire.NewSet(NewIQQDServiceImpl)
