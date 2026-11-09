package daemonizeservice

import (
	"context"
	"github.com/gin-gonic/gin"
	catalogresponse "nova-factory-server/app/business/iot/daemonize/daemonizemodels/response"
	"nova-factory-server/app/utils/gateway/v1/config/pipeline"
)

type IGatewayConfigService interface {
	Catalog(context.Context, int64) (*catalogresponse.Catalog, error)
	Generate(c *gin.Context, gatewayId int64) (*pipeline.PipelineConfig, error)
}
