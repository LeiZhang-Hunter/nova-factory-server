// 管家婆全渠道 API 的预定义错误常量。
// 覆盖 OAuth 授权、Token 签发、API 签名及数据操作中可能出现的各类错误。
package models

import (
	"errors"
)

var (
	// ErrInvalidCredential 无效的应用凭据（appKey/appSecret 不匹配）
	ErrInvalidCredential = errors.New("invalid app credential")
	// ErrInvalidRedirectURI 无效的回调地址
	ErrInvalidRedirectURI = errors.New("invalid redirect_uri")
	// ErrInvalidGrantType 不支持的 grant_type
	ErrInvalidGrantType = errors.New("invalid grant_type")
	// ErrInvalidAuthCode 无效或已过期的授权码
	ErrInvalidAuthCode = errors.New("invalid authorization code")
	// ErrInvalidRefreshToken 无效或已过期的 refresh_token
	ErrInvalidRefreshToken = errors.New("invalid refresh_token")
	// ErrInvalidAccessToken 无效或已过期的 access_token
	ErrInvalidAccessToken = errors.New("invalid access_token")
	// ErrInvalidSign 请求签名校验失败
	ErrInvalidSign = errors.New("invalid sign")
	// ErrDatabaseRequired 数据库不可用（商品/库存操作需要数据库）
	ErrDatabaseRequired = errors.New("qqd database is required")
	// ErrProductIDRequired 缺少商品ID
	ErrProductIDRequired = errors.New("productid is required")
	// ErrProductQtyRequired 未传 skus 时必须传 productqty
	ErrProductQtyRequired = errors.New("productqty is required when skus is empty")
	// ErrInvalidProductQty 无效的商品数量格式
	ErrInvalidProductQty = errors.New("invalid productqty")
	// ErrInvalidSkus 无效的 SKU 格式（应为 skuid:qty,skuid:qty 形式）
	ErrInvalidSkus = errors.New("invalid skus")
)
