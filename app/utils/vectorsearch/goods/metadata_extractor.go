package goods

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/viper"

	searchutil "nova-factory-server/app/utils/vectorsearch"
	"nova-factory-server/app/utils/vectorsearch/normalization"
	normalizationapi "nova-factory-server/app/utils/vectorsearch/normalization/api"
	"nova-factory-server/app/utils/yaml"
)

const goodsMetadataExtractorConfigKey = "vectorsearch.normalization.goods"

type metadataExtractorConfig struct {
	Normalize any            `mapstructure:"normalize"`
	Metadata  map[string]any `mapstructure:"metadata"`
}

type namedPipeline struct {
	key      string
	pipeline *normalization.Pipeline
}

// MetadataExtractor 按配置编排商品归一化与元数据提取流程。
// 当前约定：
// 1. `normalize` 是基础归一化 pipeline；
// 2. `metadata.<name>` 支持两种写法：
//   - 直接写 step 数组，例如 `spec: [...]`
//   - 写成完整 pipeline 配置，例如 `category.interceptors: [...]`
type MetadataExtractor struct {
	normalizePipeline *normalization.Pipeline
	metadataPipelines []namedPipeline
}

// NewMetadataExtractor 创建商品元数据提取器。
func NewMetadataExtractor() *MetadataExtractor {
	return &MetadataExtractor{}
}

// Init 从配置中加载并构造 normalize / metadata 对应的 pipeline。
func (e *MetadataExtractor) Init() error {
	rawConfig := metadataExtractorConfig{}
	if err := viper.UnmarshalKey(goodsMetadataExtractorConfigKey, &rawConfig); err != nil {
		return fmt.Errorf("读取商品元数据提取配置失败: %w", err)
	}

	e.normalizePipeline = nil
	e.metadataPipelines = nil

	normalizeConfig, err := decodeMetadataPipelineConfig(rawConfig.Normalize)
	if err != nil {
		return fmt.Errorf("解析商品归一化配置失败: %w", err)
	}
	if len(normalizeConfig.Interceptors) > 0 {
		e.normalizePipeline = normalization.NewPipeline(normalizeConfig)
	}

	if len(rawConfig.Metadata) == 0 {
		return nil
	}

	keys := make([]string, 0, len(rawConfig.Metadata))
	for key := range rawConfig.Metadata {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	pipelines := make([]namedPipeline, 0, len(keys))
	for _, key := range keys {
		config, err := decodeMetadataPipelineConfig(rawConfig.Metadata[key])
		if err != nil {
			return fmt.Errorf("解析商品元数据配置[%s]失败: %w", key, err)
		}
		if len(config.Interceptors) == 0 {
			continue
		}
		pipelines = append(pipelines, namedPipeline{
			key:      key,
			pipeline: normalization.NewPipeline(config),
		})
	}
	e.metadataPipelines = pipelines
	return nil
}

// Extract 执行基础归一化，并根据 metadata 配置补充规格、分类等元数据。
func (e *MetadataExtractor) Extract(input string) (normalization.Result, error) {
	result := normalization.Result{
		Original: input,
		Value:    input,
	}
	if e == nil {
		return result, nil
	}

	var err error
	if e.normalizePipeline != nil {
		result, err = e.normalizePipeline.Normalize(input)
		if err != nil {
			return normalization.Result{}, err
		}
	}
	baseValue := strings.TrimSpace(result.Value)
	if baseValue == "" {
		result.Value = ""
		return result, nil
	}

	for _, item := range e.metadataPipelines {
		if item.pipeline == nil {
			continue
		}
		extracted, runErr := item.pipeline.Normalize(baseValue)
		if runErr != nil {
			return normalization.Result{}, fmt.Errorf("执行商品元数据提取[%s]失败: %w", item.key, runErr)
		}
		mergeNormalizationResult(&result, extracted)
		for _, value := range collectMetadataValues(item.key, baseValue, extracted) {
			appendMetadataValue(&result, item.key, value)
		}
	}

	result.Value = strings.TrimSpace(result.Value)
	return result, nil
}

func decodeMetadataPipelineConfig(raw any) (normalizationapi.Config, error) {
	if raw == nil {
		return normalizationapi.Config{}, nil
	}

	content, err := yaml.Marshal(raw)
	if err != nil {
		return normalizationapi.Config{}, err
	}

	config := normalizationapi.Config{}
	if err = yaml.UnmarshalWithPrettyError(content, &config); err == nil && len(config.Interceptors) > 0 {
		return config, nil
	}

	interceptors := make([]*normalizationapi.InterceptorConfig, 0)
	if err = yaml.UnmarshalWithPrettyError(content, &interceptors); err != nil {
		return normalizationapi.Config{}, err
	}
	return normalizationapi.Config{Interceptors: interceptors}, nil
}

func mergeNormalizationResult(target *normalization.Result, incoming normalization.Result) {
	if target == nil {
		return
	}
	for _, category := range incoming.Categories {
		appendCategory(target, category)
	}
	for _, match := range incoming.Matches {
		appendMatch(target, match)
	}
	for key, values := range incoming.Metadata {
		for _, value := range values {
			appendMetadataValue(target, key, value)
		}
	}
}

func collectMetadataValues(key, baseValue string, result normalization.Result) []string {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil
	}
	if values, ok := result.Metadata[key]; ok && len(values) > 0 {
		return append([]string(nil), values...)
	}
	if key == "category" {
		return nil
	}
	if key == "spec" {
		if processed := searchutil.ProcessQuery(result.Value); processed != nil && len(processed.SpecTerms) > 0 {
			return append([]string(nil), processed.SpecTerms...)
		}
	}
	if len(result.Matches) == 0 {
		return nil
	}
	value := searchutil.NormalizeWhitespace(result.Value)
	if value == "" {
		return nil
	}
	if key != "spec" && value == searchutil.NormalizeWhitespace(baseValue) {
		return nil
	}
	return []string{value}
}

func appendCategory(result *normalization.Result, category normalizationapi.Category) {
	if result == nil || strings.TrimSpace(category.Name) == "" {
		return
	}
	for _, item := range result.Categories {
		if item.Name == category.Name && item.ID == category.ID {
			return
		}
	}
	result.Categories = append(result.Categories, category)
}

func appendMatch(result *normalization.Result, match normalizationapi.Match) {
	if result == nil {
		return
	}
	result.Matches = append(result.Matches, match)
}

func appendMetadataValue(result *normalization.Result, key, value string) {
	if result == nil {
		return
	}
	key = searchutil.NormalizeWhitespace(key)
	value = searchutil.NormalizeWhitespace(value)
	if key == "" || value == "" {
		return
	}
	if result.Metadata == nil {
		result.Metadata = make(map[string][]string)
	}
	for _, item := range result.Metadata[key] {
		if item == value {
			return
		}
	}
	result.Metadata[key] = append(result.Metadata[key], value)
}
