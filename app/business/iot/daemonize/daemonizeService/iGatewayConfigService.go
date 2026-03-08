package daemonizeService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/utils/gateway/v1/config/pipeline"
)

type IGatewayConfigService interface {
	Generate(c *gin.Context, gatewayId int64) (*pipeline.PipelineConfig, error)
}
