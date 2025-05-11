package monitorController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/monitor/monitorModels"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
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

	baizeContext.SuccessData(c, monitorModels.NewServer())
}
