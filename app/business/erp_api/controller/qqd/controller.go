package qqd

import (
	qqdservice "nova-factory-server/app/business/erp_api/service/qqd"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewAuthorize, NewToken, NewAPI, wire.Struct(new(Controller), "*"))

type Controller struct {
	Authorize *Authorize
	Token     *Token
	API       *API
}

func (q *Controller) PublicRoutes(group *gin.RouterGroup) {
	qqd := group.Group("/qqd")
	{
		qqd.GET("/oauth/authorize", q.Authorize.Authorize)
		qqd.POST("/oauth/token", q.Token.Token)
		qqd.POST("/api", q.API.API)
	}
}

type baseController struct {
	service qqdservice.Service
}
