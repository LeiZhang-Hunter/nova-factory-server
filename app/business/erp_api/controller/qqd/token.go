package qqd

import (
	"errors"
	"net/http"

	qqdservice "nova-factory-server/app/business/erp_api/service/qqd"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Token struct {
	baseController
}

func NewToken(service qqdservice.Service) *Token {
	return &Token{baseController: baseController{service: service}}
}

func (q *Token) Token(c *gin.Context) {
	response, err := q.service.IssueToken(
		c,
		c.PostForm("appkey"),
		c.PostForm("appsecret"),
		c.PostForm("code"),
		c.PostForm("grant_type"),
		c.PostForm("refresh_token"),
	)
	if err != nil {
		if !isKnownTokenError(err) {
			zap.L().Error("generate qqd token failed", zap.Error(err))
		}
		c.JSON(http.StatusOK, qqdError(tokenErrorMessage(err)))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"iserror":         false,
		"errormsg":        "ok",
		"token":           response.Token,
		"expiredate":      response.ExpireDate,
		"refresh_token":   response.RefreshToken,
		"re_expiredate":   response.RefreshExpireAt,
		"appkey":          response.AppKey,
		"appsecret":       response.AppSecret,
		"selfmallaccount": response.SelfMallAccount,
	})
}

func tokenErrorMessage(err error) string {
	if isKnownTokenError(err) {
		return err.Error()
	}
	return "generate token failed"
}

func isKnownTokenError(err error) bool {
	return errors.Is(err, qqdservice.ErrInvalidCredential) ||
		errors.Is(err, qqdservice.ErrInvalidAuthCode) ||
		errors.Is(err, qqdservice.ErrInvalidRefreshToken) ||
		errors.Is(err, qqdservice.ErrInvalidGrantType)
}
