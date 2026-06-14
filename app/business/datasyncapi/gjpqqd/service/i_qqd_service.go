// 管家婆全渠道服务接口定义。
// GjpQqdService 定义了 OAuth 授权、Token 签发、签名校验、
// 商品管理及库存更新等核心业务能力，
// 由 impl 包中的 IQQDServiceImpl 实现。
package service

import (
	"nova-factory-server/app/business/datasyncapi/gjpqqd/models"
	"nova-factory-server/app/utils/store/goods"

	"github.com/gin-gonic/gin"
)

// GjpQqdService 管家婆全渠道核心服务接口
// 提供 OAuth 授权码生成、Token 签发/刷新、访问令牌校验、
// 签名校验以及商品/库存的 CRUD 操作
type GjpQqdService interface {
	// CreateAuthorizationCallback 校验 app 凭据并生成授权码，返回拼接好 code 和 state 的回调 URL
	CreateAuthorizationCallback(ctx *gin.Context, req *models.AuthorizeReq) (string, error)
	// IssueToken 签发或刷新 access_token
	// grantType 为 authorization_code 时使用 code 换取，为 refresh_token 时使用 oldRefreshToken 刷新
	IssueToken(ctx *gin.Context, appKey, appSecret, code, grantType, oldRefreshToken string) (models.TokenResponse, error)

	// ValidAccessToken 校验 access_token 是否有效且未过期，并与 appKey 匹配
	ValidAccessToken(ctx *gin.Context, token, appKey string) bool
	// ValidSign 校验请求的 MD5 签名是否与预期一致
	ValidSign(params map[string]string, body, sign string, config *models.QQDConfig) bool
	// ProductList 分页查询商品列表，返回管家婆 API 兼容的响应格式
	ProductList(ctx *gin.Context, request *models.ProductListRequest) goods.DataResult
	// GetConfig 读取配置
	GetConfig(ctx *gin.Context) (*models.QQDConfig, error)
}
