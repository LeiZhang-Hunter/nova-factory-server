package agent

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"nova-factory-server/app/business/ai/core/client"
	"nova-factory-server/app/business/ai/core/gateway/api"
)

type Conversations struct {
	client *client.Client
}

func NewConversations(client *client.Client) api.Conversations {
	return &Conversations{
		client: client,
	}
}

func (c *Conversations) Chat(ctx context.Context, req *api.SendMessageInput) (*api.ChatResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	if req.ConversationID == 0 {
		return nil, errors.New("conversation_id不能为空")
	}
	if strings.TrimSpace(req.Content) == "" {
		return nil, errors.New("question不能为空")
	}

	chatBodyJSON, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, _, err := c.client.DoRaw(ctx, client.Request{
		Method:       "POST",
		Path:         "/api/agent/chat",
		Headers:      map[string]string{},
		AgentGateway: req.AgentGateway,
		Body:         chatBodyJSON,
		Stream:       true,
	})
	if err != nil {
		return nil, err
	}
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(strings.ToLower(contentType), "text/event-stream") {
		return &api.ChatResponse{
			StatusCode: resp.StatusCode,
			IsStream:   true,
			Headers: map[string]string{
				"Content-Type":  resp.Header.Get("Content-Type"),
				"Cache-Control": resp.Header.Get("Cache-Control"),
				"Connection":    resp.Header.Get("Connection"),
			},
			Body: resp.Body,
		}, nil
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	message := string(data)
	if message == "" {
		message = http.StatusText(resp.StatusCode)
	}
	return &api.ChatResponse{
		StatusCode: resp.StatusCode,
		Message:    message,
	}, nil
}

// StopGeneration 停止指定会话下的模型生成。
func (c *Conversations) StopGeneration(ctx context.Context, req *api.StopGenerationInput) (*api.StopGenerationResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	if req.ConversationID == 0 {
		return nil, errors.New("conversation_id不能为空")
	}
	if strings.TrimSpace(req.TabID) == "" {
		return nil, errors.New("tab_id不能为空")
	}
	bodyJSON, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	statusCode, message, err := c.client.Do(ctx, client.Request{
		Method: "POST",
		Path:   "/api/agent/chat/stop",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		AgentGateway: req.AgentGateway,
		Body:         bodyJSON,
	})
	if err != nil {
		return nil, err
	}
	return &api.StopGenerationResponse{
		StatusCode: statusCode,
		Message:    message,
	}, nil
}
