package ai

import (
	"nova-factory-server/app/business/ai/aiDataSetController"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var GinProviderSet = wire.NewSet(NewGinEngine)

func NewGinEngine(
	r *gin.Engine,
	mpcServer *gin_mcp.GinMCP,
	ai *aiDataSetController.AiDataSet) {

	group := r.Group("")
	ai.Dataset.PublicRoutes(group)
	ai.Dataset.PrivateRoutes(group)    // 工业智能体
	ai.Prediction.PrivateRoutes(group) // 工业智能体
	ai.Exception.PrivateRoutes(group)  // 工业智能体
	ai.Control.PrivateRoutes(group)
}
