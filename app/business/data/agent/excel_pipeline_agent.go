package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"nova-factory-server/app/business/data/dao"
	"nova-factory-server/app/business/data/models/entity"
	chatmodelutil "nova-factory-server/app/utils/einoAgent"
	chatmodelstore "nova-factory-server/app/utils/store/chatmodel"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

type Service struct {
	modelConfigDao dao.IModelConfigDAO
}

var ProviderSet = wire.NewSet(NewService)

func NewService(modelConfigDao dao.IModelConfigDAO) *Service {
	return &Service{modelConfigDao: modelConfigDao}
}

const (
	// defaultModelTimeout 模型请求超时，data_model_config 暂无该字段，固定默认值。
	defaultModelTimeout = 120 * time.Second
)

// Generate 每次请求时解析模型配置并初始化模型后生成 Pipeline 配置草稿。
func (s *Service) Generate(c *gin.Context, header *multipart.FileHeader, sourceType string) (json.RawMessage, error) {
	if s == nil || s.modelConfigDao == nil {
		return nil, errors.New("数据平台 Agent 未初始化")
	}
	modelCfg, err := s.resolveModelConfig(c)
	if err != nil {
		return nil, err
	}
	m, err := chatmodelutil.NewChatModel(c, modelCfg.Protocol, *modelCfg)
	if err != nil {
		return nil, fmt.Errorf("创建模型失败（protocol=%s）: %w", modelCfg.Protocol, err)
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
	out, err := chatmodelutil.Generate(c, m, "你必须输出严格 JSON 配置。", prompt)
	if err != nil {
		return nil, err
	}
	return normalizeModelJSON(out.Content)
}

// resolveModelConfig 读取 data_model_config 表，并经由 store 层获取 AI 模块连接信息后合并。
func (s *Service) resolveModelConfig(c *gin.Context) (*chatmodelutil.Config, error) {
	row, err := s.modelConfigDao.Get(c)
	if err != nil {
		return nil, err
	}
	if row == nil || strings.TrimSpace(row.Provider) == "" || strings.TrimSpace(row.Model) == "" {
		return nil, errors.New("数据平台未配置模型，请先在「配置管理 → 模型配置」中设置模型供应商与默认模型")
	}
	conn, err := chatmodelstore.GetStore().GetConnection(c, row.Provider, row.Model)
	if err != nil {
		return nil, fmt.Errorf("读取模型连接信息失败: %w", err)
	}
	if conn == nil {
		return nil, fmt.Errorf("AI 模块未配置模型 %q 的连接信息（api_type/api_key）", row.Model)
	}
	return buildChatModelConfig(row, conn)
}

// buildChatModelConfig 合并 data_model_config 与 AI 连接信息，生成模型初始化配置。
func buildChatModelConfig(row *entity.ModelConfig, conn chatmodelstore.LlmConnection) (*chatmodelutil.Config, error) {
	if row == nil || conn == nil {
		return nil, errors.New("模型配置或连接信息不能为空")
	}
	protocol := strings.ToLower(strings.TrimSpace(conn.GetAPIType()))
	if protocol == "" {
		return nil, errors.New("模型连接缺少 api_type")
	}
	// max_tokens 仅在启用时下发；未启用则传 0（适配层会省略该参数，交由模型使用默认值）。
	// 注意：不能回退使用 ai_user_llm.max_tokens，该字段是模型上下文长度（可能高达百万），
	// 远超各家模型 max_tokens 输出上限，会触发 400。
	maxTokens := 0
	if row.EnableMaxTokens && row.MaxTokens > 0 {
		maxTokens = row.MaxTokens
	}
	cfg := &chatmodelutil.Config{
		Protocol:  protocol,
		APIKey:    conn.GetAPIKey(),
		BaseURL:   strings.TrimRight(strings.TrimSpace(conn.GetAPIBase()), "/"),
		Model:     strings.TrimSpace(row.Model),
		MaxTokens: maxTokens,
		Timeout:   defaultModelTimeout,
	}
	if row.EnableTemperature {
		temperature := cast.ToFloat32(row.Temperature)
		cfg.Temperature = &temperature
	}
	if row.EnableTopP {
		topP := cast.ToFloat32(row.TopP)
		cfg.TopP = &topP
	}
	return cfg, nil
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
	return cast.ToString(b), nil
}
