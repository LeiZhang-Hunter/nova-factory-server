// 数据同步 API 模块，对外提供 ERP 数据同步的 REST 接口。
// 当前主要对接管家婆全渠道（gjpqqd），实现 OAuth 授权、
// Token 签发、商品同步及库存更新等功能。
package datasyncapi

// DataSyncApi 数据同步 API 模块的空结构体，仅用于 wire 依赖注入标识
type DataSyncApi struct{}
