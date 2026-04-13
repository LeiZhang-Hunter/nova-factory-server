package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"nova-factory-server/app/setting"
	"nova-factory-server/app/utils/aes"
	"nova-factory-server/app/utils/baizeContext"
	"strconv"
	"strings"

	"nova-factory-server/app/business/ai/core/client"
	"nova-factory-server/app/business/ai/core/gateway/api"

	"github.com/gin-gonic/gin"
)

type Conversations struct {
	client *client.Client
	aesKey string
}

func NewConversations(client *client.Client) api.Conversations {
	c := &Conversations{
		client: client,
	}

	c.aesKey = setting.Conf.AesKey

	return c
}

func (c *Conversations) Chat(ctx *gin.Context, req *api.SendMessageInput) (*api.ChatResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	if req.ConversationID == 0 {
		return nil, errors.New("conversation_id不能为空")
	}
	if strings.TrimSpace(req.Content) == "" {
		return nil, errors.New("question不能为空")
	}

	headers := map[string]string{}
	requestBody, contentType, err := buildChatRequestBody(ctx, req)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		headers["Content-Type"] = contentType
	}
	userId := strconv.FormatInt(req.UserID, 10)
	encryptString, err := aes.EncryptString([]byte(c.aesKey), userId)
	if err != nil {
		return nil, err
	}
	headers["X-User-Id"] = encryptString

	resp, _, err := c.client.DoRaw(ctx, client.Request{
		Method:  "POST",
		Path:    "/api/agent/chat",
		Headers: headers,
		Body:    requestBody,
		Stream:  true,
	})
	if err != nil {
		return nil, err
	}
	respContentType := resp.Header.Get("Content-Type")
	if strings.Contains(strings.ToLower(respContentType), "text/event-stream") {
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

func buildChatRequestBody(ctx *gin.Context, req *api.SendMessageInput) ([]byte, string, error) {
	form, err := ctx.MultipartForm()
	if err == nil && form != nil && len(form.File) > 0 {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		metadata := strings.TrimSpace(ctx.PostForm("metadata"))
		if metadata == "" {
			metadataBytes, marshalErr := json.Marshal(req)
			if marshalErr != nil {
				return nil, "", marshalErr
			}
			metadata = string(metadataBytes)
		}
		if err = writer.WriteField("metadata", metadata); err != nil {
			return nil, "", err
		}
		headers := form.File["file"]
		for _, header := range headers {
			file, openErr := header.Open()
			if openErr != nil {
				return nil, "", openErr
			}
			formFile, createErr := writer.CreateFormFile("file", header.Filename)
			if createErr != nil {
				_ = file.Close()
				return nil, "", createErr
			}
			if _, copyErr := io.Copy(formFile, file); copyErr != nil {
				_ = file.Close()
				return nil, "", copyErr
			}
			_ = file.Close()
		}
		if err = writer.Close(); err != nil {
			return nil, "", err
		}
		return body.Bytes(), writer.FormDataContentType(), nil
	}

	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, "", err
	}
	return requestBody, "application/json", nil
}

// StopGeneration 停止指定会话下的模型生成。
func (c *Conversations) StopGeneration(ctx *gin.Context, req *api.StopGenerationInput) (*api.StopGenerationResponse, error) {
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

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	userId := baizeContext.GetUserId(ctx)
	userIdStr := strconv.FormatInt(userId, 10)
	encryptString, err := aes.EncryptString([]byte(c.aesKey), userIdStr)
	if err != nil {
		return nil, err
	}
	headers["X-User-Id"] = encryptString

	statusCode, message, err := c.client.Do(ctx, client.Request{
		Method:       "POST",
		Path:         "/api/agent/chat/stop",
		Headers:      headers,
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
