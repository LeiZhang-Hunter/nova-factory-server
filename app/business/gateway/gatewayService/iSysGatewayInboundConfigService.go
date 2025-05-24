package gatewayService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/gateway/gatewayModels"
)

type ISysGatewayInboundConfigService interface {
	Add(c *gin.Context, config *gatewayModels.SysSetGatewayInboundConfig) (*gatewayModels.SysGatewayInboundConfig, error)
	Update(c *gin.Context, config *gatewayModels.SysSetGatewayInboundConfig) (*gatewayModels.SysGatewayInboundConfig, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *gatewayModels.SysSetGatewayInboundConfigReq) (*gatewayModels.SysSetGatewayInboundConfigList, error)
}
