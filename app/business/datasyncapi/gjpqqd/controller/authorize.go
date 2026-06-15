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

// Authorize OAuth 授权控制器，处理 GET /oauth/authorize 请求
type Authorize struct {
	service service.GjpQqdService
}

// NewAuthorize 创建授权控制器实例
func NewAuthorize(qqdService service.GjpQqdService) *Authorize {
	return &Authorize{
		service: qqdService,
	}
}

// Authorize 处理 OAuth 授权码请求
// - 绑定请求参数（appkey, appsecret, redirect_uri, state）
// - 校验 app 凭据有效性
// - 生成一次性授权码并写入缓存
// - 执行 302 重定向，将授权码和 state 带回 redirect_uri
// @Summary OAuth 授权码获取
// @Description 第三方系统跳转授权页面，校验 appkey/appsecret 后生成一次性授权码，302 重定向回 redirect_uri 并携带授权码和 state
// @Tags 数据同步API/管家婆全渠道
// @Param appkey query string true "应用 key"
// @Param appsecret query string true "应用 secret"
// @Param redirect_uri query string true "回调地址"
// @Param state query string false "状态参数"
// @Produce application/json
// @Success 302 "重定向到 redirect_uri 并携带授权码"
// @Failure 400 {object} map[string]interface{} "参数错误或回调地址无效"
// @Failure 401 {object} map[string]interface{} "无效凭据"
// @Failure 500 {object} map[string]interface{} "生成授权码失败"
// @Router /api/v1/erp-api/qqd/oauth/authorize [get]
func (q *Authorize) Authorize(c *gin.Context) {
	req := new(models.AuthorizeReq)
	err := c.BindQuery(req)
	if err != nil {
		zap.L().Error("binding query param failed", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	callback, err := q.service.CreateAuthorizationCallback(
		c,
		req,
	)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidCredential):
			c.JSON(http.StatusUnauthorized, qqdError(err.Error()))
		case errors.Is(err, models.ErrInvalidRedirectURI):
			c.JSON(http.StatusBadRequest, qqdError(err.Error()))
		default:
			zap.L().Error("generate qqd auth code failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, qqdError("generate auth code failed"))
		}
		return
	}
	c.Redirect(http.StatusFound, callback)
	return
}
