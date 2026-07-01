//go:build ai

package ai

import (
	"nova-factory-server/app/business/ai/agent/aidatasetcontroller"
	"nova-factory-server/app/business/ai/gateway/gatewaycontroller"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
	"go.uber.org/zap"
)

var GinProviderSet = wire.NewSet(NewFactoryBootstrap, NewGinEngine)

func NewGinEngine(
	app *routes.App,
	cache cache.Cache,
	ai *aidatasetcontroller.AiDataSet,
	gateway *gatewaycontroller.Controller,
	bootstrap *FactoryBootstrap) *AI {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				zap.L().Warn("async init llm factories panic", zap.Any("panic", r))
			}
		}()
		if err := bootstrap.Init(); err != nil {
			zap.L().Warn("init llm factories failed", zap.Error(err))
		}
	}()

	group := app.Engine.Group("")
	ai.Dataset.PublicRoutes(group)

	group.Use(middlewares.NewSessionAuthMiddlewareBuilder(cache).Build())
	{
		ai.Dataset.PrivateRoutes(group)    // 工业智能体
		ai.Prediction.PrivateRoutes(group) // 工业智能体
		ai.Exception.PrivateRoutes(group)  // 工业智能体
		ai.Control.PrivateRoutes(group)
		ai.Model.PrivateRoutes(group)
		ai.OCR.PrivateRoutes(group)
		ai.Role.PrivateRoutes(group)
		ai.Auth.PrivateRoutes(group)
		gateway.AIGateway.PrivateRoutes(group)
		gateway.Agent.PrivateRoutes(group)
		gateway.Orchestration.PrivateRoutes(group)
		gateway.SubAgent.PrivateRoutes(group)
		gateway.Conversations.PrivateRoutes(group)
		gateway.Skills.PrivateRoutes(group)
		gateway.MCPServer.PrivateRoutes(group)
		gateway.Message.PrivateRoutes(group)
		gateway.ConfigPublishHistory.PrivateRoutes(group)
		gateway.AgentConfigKey.PrivateRoutes(group)
		ai.Config.PrivateHttpRoutes(group)
	}

	// mcp
	{
		ai.Dataset.PrivateMcpRoutes(app.McpServer)    // 工业智能体
		ai.Prediction.PrivateMcpRoutes(app.McpServer) // 工业智能体
		ai.Exception.PrivateMcpRoutes(app.McpServer)  // 工业智能体
		ai.Control.PrivateMcpRoutes(app.McpServer)
		ai.Model.PrivateMcpRoutes(app.McpServer)
		ai.OCR.PrivateMcpRoutes(app.McpServer)
		ai.Role.PrivateMcpRoutes(app.McpServer)
		ai.Auth.PrivateMcpRoutes(app.McpServer)
		gateway.AIGateway.PrivateMcpRoutes(app.McpServer)
		gateway.Agent.PrivateMcpRoutes(app.McpServer)
		gateway.Orchestration.PrivateMcpRoutes(app.McpServer)
		gateway.SubAgent.PrivateMcpRoutes(app.McpServer)
		gateway.Conversations.PrivateMcpRoutes(app.McpServer)
		gateway.Skills.PrivateMcpRoutes(app.McpServer)
		gateway.MCPServer.PrivateMcpRoutes(app.McpServer)
		gateway.Message.PrivateMcpRoutes(app.McpServer)
		gateway.ConfigPublishHistory.PrivateMcpRoutes(app.McpServer)
		gateway.AgentConfigKey.PrivateMcpRoutes(app.McpServer)
	}
	return &AI{}
}
