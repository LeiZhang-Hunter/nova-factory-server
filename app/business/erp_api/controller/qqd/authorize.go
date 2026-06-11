package qqd

import (
	"errors"
	"net/http"

	qqdservice "nova-factory-server/app/business/erp_api/service/qqd"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Authorize struct {
	baseController
}

func NewAuthorize(service qqdservice.Service) *Authorize {
	return &Authorize{baseController: baseController{service: service}}
}

func (q *Authorize) Authorize(c *gin.Context) {
	callback, err := q.service.CreateAuthorizationCallback(
		c,
		c.Query("appkey"),
		c.Query("appsecret"),
		c.Query("redirect_uri"),
		c.Query("state"),
	)
	if err != nil {
		switch {
		case errors.Is(err, qqdservice.ErrInvalidCredential):
			c.JSON(http.StatusUnauthorized, qqdError(err.Error()))
		case errors.Is(err, qqdservice.ErrInvalidRedirectURI):
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
