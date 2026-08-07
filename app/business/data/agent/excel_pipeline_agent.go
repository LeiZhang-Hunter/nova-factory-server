package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
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
	"github.com/google/wire"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
	"google.golang.org/genai"
)

type Service struct {
	model model.BaseChatModel
}

var ProviderSet = wire.NewSet(NewService)

func NewService() (*Service, error) {
	if !viper.GetBool("chatwitheino.enabled") {
		return &Service{}, nil
	}
	modelName := strings.TrimSpace(viper.GetString("chatwitheino.model"))
	if modelName == "" {
		return nil, errors.New("chatwitheino.model 不能为空")
	}
	timeoutText := strings.TrimSpace(viper.GetString("chatwitheino.timeout"))
	if timeoutText == "" {
		timeoutText = "120s"
	}
	timeout, err := time.ParseDuration(timeoutText)
	if err != nil {
		return nil, fmt.Errorf("chatwitheino.timeout 配置无效: %w", err)
	}
	temperature := float32(viper.GetFloat64("chatwitheino.temperature"))
	maxTokens := viper.GetInt("chatwitheino.max_tokens")
	if maxTokens <= 0 {
		maxTokens = 8192
	}
	protocol := strings.ToLower(strings.TrimSpace(viper.GetString("chatwitheino.protocol")))
	if protocol == "" {
		protocol = "openai"
	}
	m, err := newChatModel(context.Background(), protocol, chatModelConfig{
		APIKey:      viper.GetString("chatwitheino.api_key"),
		BaseURL:     strings.TrimRight(strings.TrimSpace(viper.GetString("chatwitheino.base_url")), "/"),
		Model:       modelName,
		MaxTokens:   maxTokens,
		Temperature: temperature,
		Timeout:     timeout,
	})
	if err != nil {
		return nil, fmt.Errorf("创建 ChatWithEino 模型失败（protocol=%s）: %w", protocol, err)
	}
	return &Service{model: m}, nil
}

type chatModelConfig struct {
	APIKey      string
	BaseURL     string
	Model       string
	MaxTokens   int
	Temperature float32
	Timeout     time.Duration
}

func newChatModel(ctx context.Context, protocol string, config chatModelConfig) (model.BaseChatModel, error) {
	switch protocol {
	case "openai":
		return openaiadapter.NewChatModel(ctx, &openaiadapter.ChatModelConfig{
			APIKey:      config.APIKey,
			BaseURL:     config.BaseURL,
			Model:       config.Model,
			MaxTokens:   &config.MaxTokens,
			Temperature: &config.Temperature,
			Timeout:     config.Timeout,
		})
	case "qwen":
		return qwenadapter.NewChatModel(ctx, &qwenadapter.ChatModelConfig{
			APIKey:      config.APIKey,
			BaseURL:     config.BaseURL,
			Model:       config.Model,
			MaxTokens:   &config.MaxTokens,
			Temperature: &config.Temperature,
			Timeout:     config.Timeout,
		})
	case "ark":
		return arkadapter.NewChatModel(ctx, &arkadapter.ChatModelConfig{
			APIKey:      config.APIKey,
			BaseURL:     config.BaseURL,
			Model:       config.Model,
			MaxTokens:   &config.MaxTokens,
			Temperature: &config.Temperature,
			Timeout:     &config.Timeout,
		})
	case "ollama":
		return ollamaadapter.NewChatModel(ctx, &ollamaadapter.ChatModelConfig{
			BaseURL: config.BaseURL,
			Model:   config.Model,
			Timeout: config.Timeout,
			Options: &ollamaadapter.Options{Temperature: config.Temperature, NumPredict: config.MaxTokens},
		})
	case "gemini":
		clientConfig := &genai.ClientConfig{APIKey: config.APIKey}
		if config.BaseURL != "" {
			clientConfig.HTTPOptions = genai.HTTPOptions{BaseURL: config.BaseURL}
		}
		client, err := genai.NewClient(ctx, clientConfig)
		if err != nil {
			return nil, fmt.Errorf("创建 Gemini 客户端失败: %w", err)
		}
		return geminiadapter.NewChatModel(ctx, &geminiadapter.Config{
			Client:      client,
			Model:       config.Model,
			MaxTokens:   &config.MaxTokens,
			Temperature: &config.Temperature,
		})
	case "deepseek":
		return deepseekadapter.NewChatModel(ctx, &deepseekadapter.ChatModelConfig{
			APIKey:      config.APIKey,
			BaseURL:     config.BaseURL,
			Model:       config.Model,
			MaxTokens:   config.MaxTokens,
			Temperature: config.Temperature,
			Timeout:     config.Timeout,
		})
	case "claude":
		var baseURL *string
		if config.BaseURL != "" {
			baseURL = &config.BaseURL
		}
		return claudeadapter.NewChatModel(ctx, &claudeadapter.Config{
			APIKey:      config.APIKey,
			BaseURL:     baseURL,
			Model:       config.Model,
			MaxTokens:   config.MaxTokens,
			Temperature: &config.Temperature,
		})
	default:
		return nil, fmt.Errorf("暂不支持模型协议 %q", protocol)
	}
}

func (s *Service) Generate(ctx context.Context, header *multipart.FileHeader, sourceType string) (json.RawMessage, error) {
	if s == nil || s.model == nil {
		return nil, errors.New("ChatWithEino 未启用或模型未配置")
	}
	f, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := summarizeExcel(f)
	if err != nil {
		return nil, err
	}
	prompt := fmt.Sprintf(`你是数据采集规则专家。根据 Excel 文件结构生成一个可直接编辑的 Nova PipelineBundle 草稿。
只返回合法 JSON 对象，不要 Markdown，不要解释，不要省略必填组件。

必须遵循以下结构：
{
  "apiVersion": "nova.data/v1",
  "kind": "PipelineBundle",
  "pipelines": [{
    "name": "唯一且有意义的英文名称",
    "source": {
      "type": %q,
      "name": "file-source",
      "path": "/data/**/*.xlsx",
      "format": "excel",
      "maxBodyBytes": 1048576,
      "readHeaderTimeout": "5s",
      "metadata": {},
      "enableOsWatch": true,
      "scanTimeInterval": "10s"
    },
    "queue": {"type":"channel","name":"main-queue","capacity":1024,"batchSize":100,"batchTimeout":"1s"},
    "ontologySinks": [
      {"type":"memory","name":"local-memory"},
      {"type":"metrics","name":"metric-reporter","includeItems":true,"failOnError":false,"debug":false}
    ],
    "sourceInterceptors": [{
      "type": "parser",
      "name": "document-parser",
      "parserType": "document",
      "entityType": "Document",
      "maxRows": 10000,
      "maxTextBytes": 1048576,
      "includeRawText": false,
      "xlsxMetricRule": {
        "fileInfo": {"fileName": %q,"validSheet":[],"emptySheet":[],"mergeHeaderMode":"","globalBlockSplitRule":{"blankRowSeparateCount":2,"desc":""}},
        "matchRule": {"mode":"any","fileNameKeyword":[],"fileNamePattern":[],"pathKeyword":[],"contentKeyword":[]},
        "parseCleanRule": {
          "emptyValueReplace":{"sourceValue":["NaN"," ",""],"targetValue":""},
          "numClean":{"targetField":[],"dataType":"number","errorValue":""},
          "textField":{"targetField":[],"dataType":"string","keepRaw":true,"lineBreakReserve":true}
        },
        "scalarMapping": [],
        "fieldMapping": [],
        "computeDerivedIndex": []
      }
    }],
    "sinkInterceptors": []
  }]
}

请根据 Sheet 名、标题、表头和样例值填充 validSheet、matchRule、scalarMapping、fieldMapping 和派生指标。字段名使用清晰的 lowerCamelCase；无法可靠判断时保留空数组，不要虚构业务含义。

Excel 结构摘要：
%s`, sourceType, header.Filename, content)
	out, err := s.model.Generate(ctx, []*schema.Message{{Role: schema.System, Content: "你必须输出严格 JSON 配置。"}, {Role: schema.User, Content: prompt}})
	if err != nil {
		return nil, err
	}
	return normalizeModelJSON(out.Content)
}

func normalizeModelJSON(content string) (json.RawMessage, error) {
	text := strings.TrimSpace(content)
	if strings.HasPrefix(text, "```") {
		if lineEnd := strings.IndexByte(text, '\n'); lineEnd >= 0 {
			text = text[lineEnd+1:]
		}
		text = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(text), "```"))
	}
	var value any
	if err := json.Unmarshal([]byte(text), &value); err != nil {
		return nil, fmt.Errorf("Agent 返回的 Pipeline 不是合法 JSON: %w", err)
	}
	return json.Marshal(value)
}

func summarizeExcel(r io.Reader) (string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return "", err
	}
	defer f.Close()
	type sheet struct {
		Name string     `json:"name"`
		Rows [][]string `json:"rows"`
	}
	result := make([]sheet, 0)
	for _, name := range f.GetSheetList() {
		rows, err := f.GetRows(name)
		if err != nil {
			return "", err
		}
		if len(rows) > 12 {
			rows = rows[:12]
		}
		if len(rows) > 0 {
			result = append(result, sheet{Name: name, Rows: rows})
		}
	}
	b, _ := json.Marshal(result)
	return string(b), nil
}
