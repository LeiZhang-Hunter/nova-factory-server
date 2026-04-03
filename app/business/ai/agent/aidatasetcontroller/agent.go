package aidatasetcontroller

import (
	"net/http"
	"time"

	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	sseUtils "nova-factory-server/app/utils/sse"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Agent struct {
	service        aidatasetservice.IAiConversationService
	gatewayService aidatasetservice.IAIGatewayService
}

// NewAgent 会话控制器构造函数。
func NewAgent(service aidatasetservice.IAiConversationService, gatewayService aidatasetservice.IAIGatewayService) *Agent {
	return &Agent{
		service:        service,
		gatewayService: gatewayService,
	}
}

// PrivateRoutes 注册会话相关私有路由。
func (agent *Agent) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/agent/conversations")
	group.GET("/list", middlewares.HasPermission("ai:agent:conversations:list"), agent.ListConversations)
	group.POST("/create", middlewares.HasPermission("ai:agent:conversations:create"), agent.CreateConversation)
	group.DELETE("/remove/:ids", middlewares.HasPermission("ai:agent:conversations:remove"), agent.RemoveConversation)
	group.POST("/chat", middlewares.HasPermission("ai:agent:conversations:chat"), agent.Chat)
	group.POST("/stop-generation", middlewares.HasPermission("ai:agent:conversations:stop-generation"), agent.StopGeneration)
}

// ListConversations 查询会话列表
// @Summary 查询会话列表
// @Description 查询会话列表
// @Tags 工业智能体/会话管理
// @Param object query aidatasetmodels.AiConversationQuery true "会话列表查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/agent/conversations/list [get]
func (agent *Agent) ListConversations(c *gin.Context) {
	req := new(aidatasetmodels.AiConversationQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := agent.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// CreateConversation 创建会话
// @Summary 创建会话
// @Description 创建会话
// @Tags 工业智能体/会话管理
// @Param object body aidatasetmodels.SetAiConversation true "创建会话参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "创建成功"
// @Router /ai/agent/conversations/create [post]
func (agent *Agent) CreateConversation(c *gin.Context) {
	req := new(aidatasetmodels.SetAiConversation)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := agent.service.Create(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// RemoveConversation 删除会话（软删除）
// @Summary 删除会话
// @Description 删除会话（软删除）
// @Tags 工业智能体/会话管理
// @Param ids path string true "会话ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /ai/agent/conversations/remove/{ids} [delete]
func (agent *Agent) RemoveConversation(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := agent.service.Remove(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// Chat 会话聊天
// @Summary 会话聊天
// @Description 会话聊天
// @Tags 工业智能体/会话管理
// @Param object body aidatasetmodels.SendMessageInput true "发送消息参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "发送成功"
// @Router /ai/agent/conversations/chat [post]
func (agent *Agent) Chat(c *gin.Context) {
	req := new(aidatasetmodels.SendMessageInput)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		zap.L().Error("param error", zap.Error(err))
		return
	}
	if req.Content == "" {
		baizeContext.Waring(c, "提问内容不能为空")
		return
	}
	data, err := agent.gatewayService.Chat(c, req)
	if err != nil {
		zap.L().Error("chat error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	if data != nil && data.IsStream && data.Body != nil {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		// 禁用nginx缓存,防止nginx会缓存数据导致数据流是一段一段的
		c.Writer.Header().Set("X-Accel-Buffering", "no")
		defer data.Body.Close()
		sseUtils.ApplyHeaders(c.Writer.Header(), data.Headers)
		if data.StatusCode > 0 {
			c.Status(data.StatusCode)
		} else {
			c.Status(http.StatusOK)
		}
		if streamErr := sseUtils.Transport(c.Writer, data.Body, 60*time.Second); streamErr != nil {
			if writeErr := sseUtils.WriteErrorEvent(c.Writer, streamErr.Error()); writeErr != nil {
				zap.L().Warn("sse write error event failed", zap.Error(writeErr))
			}
			zap.L().Warn("sse stream read failed", zap.Error(streamErr))
		}
		return
	}
	baizeContext.SuccessData(c, data)
}

// StopGeneration 停止会话生成
// @Summary 停止会话生成
// @Description 停止上游会话生成任务
// @Tags 工业智能体/会话管理
// @Param object body aidatasetmodels.StopGenerationInput true "停止生成参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "停止成功"
// @Router /ai/agent/conversations/stop-generation [post]
func (agent *Agent) StopGeneration(c *gin.Context) {
	req := new(aidatasetmodels.StopGenerationInput)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		zap.L().Error("param error", zap.Error(err))
		return
	}
	data, err := agent.gatewayService.StopGeneration(c, req)
	if err != nil {
		zap.L().Error("stop generation error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}
