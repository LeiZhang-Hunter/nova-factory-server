package impl

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"nova-factory-server/app/business/data/models/entity"

	"gopkg.in/yaml.v3"
)

// dispatchBundle 仅解析 PipelineBundle 的 pipelines 字段，保留原始字节。
type dispatchBundle struct {
	Pipelines []json.RawMessage `json:"pipelines"`
}

// buildMergedConfig 将多条规则的 pipelines 合并为采集器使用的 YAML 配置并计算 MD5。
// 输出顶层为 pipelines: 的 YAML（与 ontology-pipelines.yml 结构一致），
// 排序、重名校验与 MD5 在预览、日志、拉取三处共用，保证数值一致。
func buildMergedConfig(rules []entity.PipelineRule) ([]byte, string, error) {
	sorted := make([]entity.PipelineRule, len(rules))
	copy(sorted, rules)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].CreatedAt.Equal(sorted[j].CreatedAt) {
			return sorted[i].ID < sorted[j].ID
		}
		return sorted[i].CreatedAt.Before(sorted[j].CreatedAt)
	})

	merged := make([]json.RawMessage, 0, len(sorted))
	owners := make(map[string]string, 8)
	for i := range sorted {
		rule := &sorted[i]
		var bundle dispatchBundle
		if err := json.Unmarshal([]byte(rule.Config), &bundle); err != nil {
			return nil, "", fmt.Errorf("规则 %s 配置解析失败: %w", rule.Name, err)
		}
		for _, raw := range bundle.Pipelines {
			var meta struct {
				Name string `json:"name"`
			}
			if err := json.Unmarshal(raw, &meta); err != nil {
				return nil, "", fmt.Errorf("规则 %s 包含无法解析的 pipeline: %w", rule.Name, err)
			}
			name := strings.TrimSpace(meta.Name)
			if name == "" {
				return nil, "", fmt.Errorf("规则 %s 包含无名称的 pipeline", rule.Name)
			}
			if owner, duplicated := owners[name]; duplicated {
				return nil, "", fmt.Errorf("Pipeline 名称 %q 在规则 %s 与 %s 中重复", name, owner, rule.Name)
			}
			owners[name] = rule.Name
			merged = append(merged, raw)
		}
	}

	var pipelines yaml.Node
	if len(merged) > 0 {
		raw, err := json.Marshal(merged)
		if err != nil {
			return nil, "", err
		}
		if err := yaml.Unmarshal(raw, &pipelines); err != nil {
			return nil, "", fmt.Errorf("合并配置转 YAML 失败: %w", err)
		}
		if pipelines.Kind == yaml.DocumentNode && len(pipelines.Content) > 0 {
			pipelines = *pipelines.Content[0]
		}
		clearYAMLStyle(&pipelines)
	} else {
		pipelines = yaml.Node{Kind: yaml.SequenceNode}
	}
	doc := yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{{
			Kind: yaml.MappingNode,
			Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Value: "pipelines"},
				&pipelines,
			},
		}},
	}
	payload, err := yaml.Marshal(&doc)
	if err != nil {
		return nil, "", err
	}
	return payload, checksumMD5(payload), nil
}

// clearYAMLStyle 递归清除节点样式，避免 JSON 输入的 flow 风格被原样输出
func clearYAMLStyle(n *yaml.Node) {
	n.Style = 0
	for _, child := range n.Content {
		clearYAMLStyle(child)
	}
}

func checksumMD5(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}
