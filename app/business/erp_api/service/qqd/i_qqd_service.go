package qqd

import (
	"errors"

	"nova-factory-server/app/business/erp_api/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateAuthorizationCallback(ctx *gin.Context, appKey, appSecret, redirectURI, state string) (string, error)
	IssueToken(ctx *gin.Context, appKey, appSecret, code, grantType, oldRefreshToken string) (TokenResponse, error)
	ValidAccessToken(ctx *gin.Context, token, appKey string) bool
	ValidSign(params map[string]string, body, sign string) bool
	ProductList(ctx *gin.Context, request ProductListRequest) (map[string]any, error)
	AddProducts(ctx *gin.Context, goodsInfos []map[string]any) ([]map[string]any, error)
	ProductStockUpdate(ctx *gin.Context, request ProductStockUpdateRequest) (map[string]any, error)
}

type TokenResponse struct {
	Token           string
	ExpireDate      string
	RefreshToken    string
	RefreshExpireAt string
	AppKey          string
	AppSecret       string
	SelfMallAccount string
}

type ProductListRequest struct {
	PageNo   int
	PageSize int
}

type ProductStockUpdateRequest struct {
	ProductID  string
	ProductQty string
	Skus       string
}

var (
	ErrInvalidCredential   = errors.New("invalid app credential")
	ErrInvalidRedirectURI  = errors.New("invalid redirect_uri")
	ErrInvalidGrantType    = errors.New("invalid grant_type")
	ErrInvalidAuthCode     = errors.New("invalid authorization code")
	ErrInvalidRefreshToken = errors.New("invalid refresh_token")
	ErrInvalidAccessToken  = errors.New("invalid access_token")
	ErrInvalidSign         = errors.New("invalid sign")
	ErrDatabaseRequired    = errors.New("qqd database is required")
	ErrProductIDRequired   = errors.New("productid is required")
	ErrProductQtyRequired  = errors.New("productqty is required when skus is empty")
	ErrInvalidProductQty   = errors.New("invalid productqty")
	ErrInvalidSkus         = errors.New("invalid skus")
)

type QQDConfig = models.QQDConfig
