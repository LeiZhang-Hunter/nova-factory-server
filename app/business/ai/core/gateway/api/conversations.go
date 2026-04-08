package api

import (
	"context"
	"io"

	"github.com/gin-gonic/gin"
)

type Conversations interface {
	Chat(ctx context.Context, req *SendMessageInput) (*ChatResponse, error)
	StopGeneration(ctx *gin.Context, req *StopGenerationInput) (*StopGenerationResponse, error)
}

// SendMessageInput input for sending a message
type SendMessageInput struct {
	ConversationID int64  `json:"conversation_id,string"`
	AgentGateway   string `json:"agent_gateway"` // AgentGateway 指定网关标识。
	Content        string `json:"content"`
	TabID          string `json:"tab_id"`
	UserID         int64  `json:"user_id"`
}

// StopGenerationInput 停止生成请求参数。
type StopGenerationInput struct {
	ConversationID int64  `json:"conversation_id,string"` // ConversationID 会话ID。
	AgentGateway   string `json:"agent_gateway"`          // AgentGateway 指定网关标识。
	TabID          string `json:"tab_id"`                 // TabID 标签页ID。
}

// ChatRequest 网关聊天请求参数。
type ChatRequest struct {
	ChatID       string            `json:"chat_id"`       // ChatID 会话所属聊天ID。
	SessionID    string            `json:"session_id"`    // SessionID 当前会话ID。
	UserID       string            `json:"user_id"`       // UserID 当前请求用户ID。
	Question     string            `json:"question"`      // Question 用户输入问题内容。
	AgentGateway string            `json:"agent_gateway"` // AgentGateway 指定网关标识。
	Headers      map[string]string `json:"headers"`       // Headers 透传到网关的请求头。
	Stream       bool              `json:"stream"`        // Stream 是否启用流式响应。
}

type ChatResponse struct {
	StatusCode int               `json:"status_code"`
	Message    string            `json:"message"`
	IsStream   bool              `json:"is_stream"`
	Headers    map[string]string `json:"-"`
	Body       io.ReadCloser     `json:"-"`
}

// StopGenerationResponse 停止生成响应。
type StopGenerationResponse struct {
	StatusCode int    `json:"status_code"` // StatusCode 上游返回状态码。
	Message    string `json:"message"`     // Message 上游返回消息。
}
