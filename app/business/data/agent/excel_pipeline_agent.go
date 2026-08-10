package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"nova-factory-server/app/constant/aiagent"
	chatmodelutil "nova-factory-server/app/utils/einoAgent"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

type Service struct {
}

var ProviderSet = wire.NewSet(NewService)

func NewService() *Service {
	return &Service{}
}

// Generate 每次请求时解析模型配置并初始化模型后生成 Pipeline 配置草稿。
func (s *Service) Generate(c *gin.Context, header *multipart.FileHeader, sourceType string) (json.RawMessage, error) {
	if s == nil {
		return nil, errors.New("数据平台 Agent 未初始化")
	}
	modelCfg, err := chatmodelutil.BuildAgentChatModelConfig(c, aiagent.DataAgentType)
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
