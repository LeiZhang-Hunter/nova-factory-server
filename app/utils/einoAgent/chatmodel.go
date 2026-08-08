package chatmodel

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	arkadapter "github.com/cloudwego/eino-ext/components/model/ark"
	claudeadapter "github.com/cloudwego/eino-ext/components/model/claude"
	deepseekadapter "github.com/cloudwego/eino-ext/components/model/deepseek"
	geminiadapter "github.com/cloudwego/eino-ext/components/model/gemini"
	ollamaadapter "github.com/cloudwego/eino-ext/components/model/ollama"
	openaiadapter "github.com/cloudwego/eino-ext/components/model/openai"
	qwenadapter "github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/spf13/cast"
	"google.golang.org/genai"
)

// Config 初始化 eino ChatModel 所需的全部参数。
// Temperature/TopP 为 nil 表示不传该采样参数，交由模型使用默认行为；
// 值类型适配器（deepseek/ollama）在 nil 时会传 0。
type Config struct {
	Protocol    string
	APIKey      string
	BaseURL     string
	Model       string
	MaxTokens   int
	Temperature *float32
	TopP        *float32
	Timeout     time.Duration
}

// NewChatModel 根据协议创建对应的 eino ChatModel。
// 只根据入参初始化，不读取任何外部配置。
func NewChatModel(ctx context.Context, protocol string, cfg Config) (model.BaseChatModel, error) {
	protocol = strings.ToLower(strings.TrimSpace(protocol))
	if protocol == "" {
		return nil, fmt.Errorf("模型协议不能为空")
	}
	if strings.TrimSpace(cfg.Model) == "" {
		return nil, fmt.Errorf("模型名称不能为空")
	}

	// 统一换算：指针型适配器直接使用 cfg.Temperature/cfg.TopP/maxTokens，
	// 值型适配器（deepseek/ollama）使用换算后的零值兜底。
	var maxTokens *int
	if cfg.MaxTokens > 0 {
		v := cfg.MaxTokens
		maxTokens = &v
	}
	// cast.ToFloat32 会解指针，nil 转换为 0。
	temperature := cast.ToFloat32(cfg.Temperature)
	topP := cast.ToFloat32(cfg.TopP)

	switch protocol {
	case "openai":
		return openaiadapter.NewChatModel(ctx, &openaiadapter.ChatModelConfig{
			APIKey:      cfg.APIKey,
			BaseURL:     cfg.BaseURL,
			Model:       cfg.Model,
			MaxTokens:   maxTokens,
			Temperature: cfg.Temperature,
			TopP:        cfg.TopP,
			Timeout:     cfg.Timeout,
		})
	case "qwen":
		return qwenadapter.NewChatModel(ctx, &qwenadapter.ChatModelConfig{
			APIKey:      cfg.APIKey,
			BaseURL:     cfg.BaseURL,
			Model:       cfg.Model,
			MaxTokens:   maxTokens,
			Temperature: cfg.Temperature,
			TopP:        cfg.TopP,
			Timeout:     cfg.Timeout,
		})
	case "ark":
		return arkadapter.NewChatModel(ctx, &arkadapter.ChatModelConfig{
			APIKey:      cfg.APIKey,
			BaseURL:     cfg.BaseURL,
			Model:       cfg.Model,
			MaxTokens:   maxTokens,
			Temperature: cfg.Temperature,
			TopP:        cfg.TopP,
			Timeout:     &cfg.Timeout,
		})
	case "ollama":
		return ollamaadapter.NewChatModel(ctx, &ollamaadapter.ChatModelConfig{
			BaseURL: cfg.BaseURL,
			Model:   cfg.Model,
			Timeout: cfg.Timeout,
			Options: &ollamaadapter.Options{
				Temperature: temperature,
				TopP:        topP,
				NumPredict:  cfg.MaxTokens,
			},
		})
	case "gemini":
		clientConfig := &genai.ClientConfig{APIKey: cfg.APIKey}
		if cfg.BaseURL != "" {
			clientConfig.HTTPOptions = genai.HTTPOptions{BaseURL: cfg.BaseURL}
		}
		client, err := genai.NewClient(ctx, clientConfig)
		if err != nil {
			return nil, fmt.Errorf("创建 Gemini 客户端失败: %w", err)
		}
		return geminiadapter.NewChatModel(ctx, &geminiadapter.Config{
			Client:      client,
			Model:       cfg.Model,
			MaxTokens:   maxTokens,
			Temperature: cfg.Temperature,
			TopP:        cfg.TopP,
		})
	case "deepseek":
		return deepseekadapter.NewChatModel(ctx, &deepseekadapter.ChatModelConfig{
			APIKey:      cfg.APIKey,
			BaseURL:     cfg.BaseURL,
			Model:       cfg.Model,
			MaxTokens:   cfg.MaxTokens,
			Temperature: temperature,
			TopP:        topP,
			Timeout:     cfg.Timeout,
		})
	case "claude":
		var baseURL *string
		if cfg.BaseURL != "" {
			baseURL = &cfg.BaseURL
		}
		return claudeadapter.NewChatModel(ctx, &claudeadapter.Config{
			APIKey:      cfg.APIKey,
			BaseURL:     baseURL,
			Model:       cfg.Model,
			MaxTokens:   cfg.MaxTokens,
			Temperature: cfg.Temperature,
			TopP:        cfg.TopP,
		})
	default:
		return nil, fmt.Errorf("暂不支持模型协议 %q", protocol)
	}
}

// GenerateMessages 以原始消息列表调用 ChatModel 生成回复。
func GenerateMessages(ctx context.Context, m model.BaseChatModel, messages []*schema.Message) (*schema.Message, error) {
	if m == nil {
		return nil, errors.New("chat model 未初始化")
	}
	return m.Generate(ctx, messages)
}

// Generate 以 system + user 两条消息调用 ChatModel 生成回复。
func Generate(ctx context.Context, m model.BaseChatModel, system, user string) (*schema.Message, error) {
	if m == nil {
		return nil, errors.New("chat model 未初始化")
	}
	return GenerateMessages(ctx, m, []*schema.Message{
		{Role: schema.System, Content: system},
		{Role: schema.User, Content: user},
	})
}
