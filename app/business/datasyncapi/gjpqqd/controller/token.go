// 管家婆全渠道 Token 签发控制器。
// 处理第三方的 Token 请求（authorization_code / refresh_token 两种 grant_type），
// 校验授权码或 refresh_token 后签发 access_token。
package controller

import (
	"errors"
	"net/http"
	"nova-factory-server/app/business/datasyncapi/gjpqqd/models"
	"nova-factory-server/app/business/datasyncapi/gjpqqd/service"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Token Token 签发控制器，处理 POST /oauth/token 请求
type Token struct {
	gjpService service.GjpQqdService
}

// NewToken 创建 Token 控制器实例
func NewToken(gjpService service.GjpQqdService) *Token {
	return &Token{
		gjpService: gjpService,
	}
}

// Token 处理 Token 签发请求
// 支持两种 grant_type：
// - authorization_code: 使用 OAuth 授权码换取 access_token
// - refresh_token: 使用 refresh_token 刷新 access_token
// 校验失败时返回已知错误码，未知错误统一返回 "generate token failed"
func (q *Token) Token(c *gin.Context) {
	req := new(models.IssueTokenReq)
	err := c.ShouldBind(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	response, err := q.gjpService.IssueToken(
		c,
		req.AppKey,
		req.AppSecret,
		req.Code,
		req.GrantType,
		req.RefreshToken,
	)
	if err != nil {
		if !isKnownTokenError(err) {
			zap.L().Error("generate qqd token failed", zap.Error(err))
		}
		c.JSON(http.StatusOK, qqdError(tokenErrorMessage(err)))
		return
	}

	var responseIssueToken models.IssueTokenResponse
	responseIssueToken.IsError = false
	responseIssueToken.ErrorMsg = "ok"
	responseIssueToken.Token = response.Token
	responseIssueToken.ExpireDate = response.ExpireDate
	responseIssueToken.RefreshToken = response.RefreshToken
	responseIssueToken.ReExpireDate = response.ExpireDate
	responseIssueToken.AppKey = response.AppKey
	responseIssueToken.AppSecret = response.AppSecret
	responseIssueToken.SelfMallAccount = response.SelfMallAccount
	c.JSON(http.StatusOK, responseIssueToken)
}

// tokenErrorMessage 根据错误类型返回对外消息，已知错误原样返回，未知错误统一提示
func tokenErrorMessage(err error) string {
	if isKnownTokenError(err) {
		return err.Error()
	}
	return "generate token failed"
}

// isKnownTokenError 判断是否为预定义的 Token 签发错误
func isKnownTokenError(err error) bool {
	return errors.Is(err, models.ErrInvalidCredential) ||
		errors.Is(err, models.ErrInvalidAuthCode) ||
		errors.Is(err, models.ErrInvalidRefreshToken) ||
		errors.Is(err, models.ErrInvalidGrantType)
}
