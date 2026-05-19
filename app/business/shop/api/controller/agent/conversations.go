//go:build ai

package agent

import (
	"encoding/json"
	"net/http"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/sse"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Conversations 商城智能体会话控制器。
type Conversations struct {
	service        gatewayservice.IAiConversationService
	gatewayService aidatasetservice.IAIGatewayService
	agentService   gatewayservice.IAIAgentService
}

// NewConversations 创建商城智能体会话控制器。
func NewConversations(service gatewayservice.IAiConversationService, gatewayService aidatasetservice.IAIGatewayService, agentService gatewayservice.IAIAgentService) *Conversations {
	return &Conversations{
		service:        service,
		gatewayService: gatewayService,
		agentService:   agentService,
	}
}

// PrivateRoutes 注册商城智能体会话路由。
func (conversations *Conversations) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/agent/conversations")
	group.GET("/list", conversations.ListConversations)
	group.POST("/create", conversations.CreateConversation)
	group.DELETE("/remove/:ids", conversations.RemoveConversation)
	group.POST("/chat", conversations.Chat)
	group.POST("/stop-generation", conversations.StopGeneration)
}

// ListConversations 查询会话列表
// @Summary 查询会话列表
// @Description 查询当前登录用户的智能体会话列表
// @Tags app接口/商城/App智能体会话
// @Param object query aidatasetmodels.AiConversationQuery true "会话列表查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/agent/conversations/list [get]
func (conversations *Conversations) ListConversations(c *gin.Context) {
	req := new(aidatasetmodels.AiConversationQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := conversations.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// CreateConversation 创建会话
// @Summary 创建会话
// @Description 创建当前登录用户的智能体会话
// @Tags app接口/商城/App智能体会话
// @Param object body aidatasetmodels.SetAiConversation true "创建会话参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "创建成功"
// @Router /api/v1/app/shop/agent/conversations/create [post]
func (conversations *Conversations) CreateConversation(c *gin.Context) {
	req := new(aidatasetmodels.SetAiConversation)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		baizeContext.Waring(c, "名称为空")
		return
	}
	if req.AgentID == 0 {
		baizeContext.Waring(c, "智能体ID不能为空")
		return
	}

	agentInfo, err := conversations.agentService.GetByID(c, req.AgentID)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	if agentInfo == nil {
		baizeContext.Waring(c, "智能体配置不存在")
		return
	}
	if agentInfo.Enable == nil || !*agentInfo.Enable {
		baizeContext.Waring(c, "智能体配置未启用")
		return
	}

	req.AgentID = agentInfo.ID
	req.AgentType = agentInfo.Type
	req.LLMProviderID = agentInfo.DefaultLLMProviderID
	req.LLMModelID = agentInfo.DefaultLLMModelID

	data, err := conversations.service.Create(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// RemoveConversation 删除会话（软删除）
// @Summary 删除会话
// @Description 删除当前登录用户的会话（软删除）
// @Tags app接口/商城/App智能体会话
// @Param ids path string true "会话ID，多个用逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /api/v1/app/shop/agent/conversations/remove/{ids} [delete]
func (conversations *Conversations) RemoveConversation(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := conversations.service.Remove(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// Chat 会话聊天
// @Summary 会话聊天
// @Description 发送会话消息并支持流式返回
// @Tags app接口/商城/App智能体会话
// @Param metadata formData string false "会话消息元数据(JSON字符串)"
// @Param file formData file false "上传文件"
// @Param object body aidatasetmodels.SendMessageInput false "发送消息参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "发送成功"
// @Router /api/v1/app/shop/agent/conversations/chat [post]
func (conversations *Conversations) Chat(c *gin.Context) {
	req := new(aidatasetmodels.SendMessageInput)

	contentType := c.ContentType()
	if strings.HasPrefix(contentType, "multipart/form-data") || strings.TrimSpace(c.PostForm("metadata")) != "" {
		metadata := c.PostForm("metadata")
		if err := json.Unmarshal([]byte(metadata), req); err != nil {
			baizeContext.ParameterError(c)
			zap.L().Error("metadata param error", zap.Error(err))
			return
		}
	} else {
		if err := c.ShouldBindJSON(req); err != nil {
			baizeContext.ParameterError(c)
			zap.L().Error("bind json error", zap.Error(err))
			return
		}
	}

	if req.ConversationID == 0 {
		baizeContext.Waring(c, "会话id不能为空")
		return
	}
	if strings.TrimSpace(req.Content) == "" {
		baizeContext.Waring(c, "提问内容不能为空")
		return
	}
	if strings.TrimSpace(req.TabID) == "" {
		req.TabID = "team"
	}

	data, err := conversations.gatewayService.Chat(c, req)
	if err != nil {
		zap.L().Error("chat error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	if data == nil {
		baizeContext.Waring(c, "聊天服务无响应")
		return
	}
	if data.StatusCode != http.StatusOK {
		baizeContext.Waring(c, data.Message)
		return
	}
	if data.IsStream && data.Body != nil {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("X-Accel-Buffering", "no")
		defer data.Body.Close()
		sse.ApplyHeaders(c.Writer.Header(), data.Headers)
		if data.StatusCode > 0 {
			c.Writer.Header().Set("X-Upstream-Status-Code", strconv.Itoa(data.StatusCode))
		}
		c.Status(http.StatusOK)
		streamErr := sse.Transport(c.Writer, data.Body, 300*time.Second)
		if streamErr != nil {
			if writeErr := sse.WriteErrorEvent(c.Writer, streamErr.Error()); writeErr != nil {
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
// @Description 停止当前会话正在进行中的模型生成任务
// @Tags app接口/商城/App智能体会话
// @Param object body aidatasetmodels.StopGenerationInput true "停止生成参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "停止成功"
// @Router /api/v1/app/shop/agent/conversations/stop-generation [post]
func (conversations *Conversations) StopGeneration(c *gin.Context) {
	req := new(aidatasetmodels.StopGenerationInput)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		zap.L().Error("param error", zap.Error(err))
		return
	}
	if strings.TrimSpace(req.TabID) == "" {
		req.TabID = "team"
	}

	data, err := conversations.gatewayService.StopGeneration(c, req)
	if err != nil {
		zap.L().Error("stop generation error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	if data == nil {
		baizeContext.Waring(c, "停止生成服务无响应")
		return
	}
	if data.StatusCode != http.StatusOK {
		baizeContext.Waring(c, data.Message)
		return
	}
	baizeContext.SuccessData(c, data)
}
