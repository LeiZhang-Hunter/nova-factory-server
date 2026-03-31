package monitorcontroller

import (
	"nova-factory-server/app/business/admin/monitor/monitormodels"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type InfoServer struct {
}

func NewInfoServer() *InfoServer {
	return &InfoServer{}
}

func (isc *InfoServer) PrivateRoutes(router *gin.RouterGroup) {
	server := router.Group("/monitor/server")
	server.GET("", middlewares.HasPermission("monitor:server"), isc.GetInfoServer)
}

func (isc *InfoServer) GetInfoServer(c *gin.Context) {

	baizeContext.SuccessData(c, monitormodels.NewServer())
}
