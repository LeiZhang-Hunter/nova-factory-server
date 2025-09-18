package aiDataSetController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/middlewares"
)

type Exception struct{}

func NewException() *Exception {
	return &Exception{}
}

func (e *Exception) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/exception")
	group.POST("/set", middlewares.HasPermission("ai:exception:set"), e.Set)
	group.POST("/info", middlewares.HasPermission("ai:exception:set"), e.Info)
	group.DELETE("/remove/:ids", middlewares.HasPermission("ai:exception:remove"), e.Remove)
	return
}

func (e *Exception) Set(c *gin.Context) {

}

func (e *Exception) Info(c *gin.Context) {

}

func (e *Exception) Remove(c *gin.Context) {

}
