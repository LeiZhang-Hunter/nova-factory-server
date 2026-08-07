package impl

import (
	"encoding/json"
	"fmt"
	"strings"

	"nova-factory-server/app/business/data/models/dto"
	"nova-factory-server/app/business/data/models/entity"
)

func validatePipelineBundle(raw json.RawMessage, sourceType string) (*dto.PipelineBundle, []string, error) {
	var bundle dto.PipelineBundle
	if err := json.Unmarshal(raw, &bundle); err != nil {
		return nil, nil, fmt.Errorf("配置不是有效的 JSON: %w", err)
	}
	if bundle.APIVersion != "nova.data/v1" {
		return nil, nil, fmt.Errorf("apiVersion 必须为 nova.data/v1")
	}
	if bundle.Kind != "PipelineBundle" {
		return nil, nil, fmt.Errorf("kind 必须为 PipelineBundle")
	}
	if len(bundle.Pipelines) == 0 {
		return nil, nil, fmt.Errorf("至少需要一个 Pipeline")
	}

	names := make(map[string]struct{}, len(bundle.Pipelines))
	refs := make([]string, 0)
	for i, pipeline := range bundle.Pipelines {
		prefix := fmt.Sprintf("pipelines[%d]", i)
		name := strings.TrimSpace(pipeline.Name)
		if name == "" {
			return nil, nil, fmt.Errorf("%s.name 不能为空", prefix)
		}
		if _, ok := names[name]; ok {
			return nil, nil, fmt.Errorf("Pipeline 名称 %q 重复", name)
		}
		names[name] = struct{}{}

		actualType := componentString(pipeline.Source, "type")
		if actualType != sourceType {
			return nil, nil, fmt.Errorf("%s.source.type 必须为 %s", prefix, sourceType)
		}
		if componentString(pipeline.Queue, "type") != "channel" {
			return nil, nil, fmt.Errorf("%s.queue.type 目前只支持 channel", prefix)
		}
		if componentString(pipeline.Queue, "name") == "" {
			return nil, nil, fmt.Errorf("%s.queue.name 不能为空", prefix)
		}
		if len(pipeline.OntologySinks) == 0 {
			return nil, nil, fmt.Errorf("%s.ontologySinks 至少需要一项", prefix)
		}
		for j, sink := range pipeline.OntologySinks {
			t := componentString(sink, "type")
			if t != "memory" && t != "metrics" {
				return nil, nil, fmt.Errorf("%s.ontologySinks[%d].type 不支持 %q", prefix, j, t)
			}
			if componentString(sink, "name") == "" {
				return nil, nil, fmt.Errorf("%s.ontologySinks[%d].name 不能为空", prefix, j)
			}
		}
		for j, interceptor := range append(append([]dto.Component{}, pipeline.SourceInterceptors...), pipeline.SinkInterceptors...) {
			t := componentString(interceptor, "type")
			if t != "parser" && t != "mapper" {
				return nil, nil, fmt.Errorf("%s.interceptors[%d].type 不支持 %q", prefix, j, t)
			}
		}

		switch sourceType {
		case entity.SourceTypeFile:
			if componentString(pipeline.Source, "path") == "" {
				return nil, nil, fmt.Errorf("%s.source.path 不能为空", prefix)
			}
			if err := validateExcelParser(pipeline.SourceInterceptors, prefix); err != nil {
				return nil, nil, err
			}
		case entity.SourceTypeMySQL:
			ref := componentString(pipeline.Source, "connectionRef")
			if ref == "" {
				return nil, nil, fmt.Errorf("%s.source.connectionRef 不能为空", prefix)
			}
			if componentString(pipeline.Source, "query") == "" {
				return nil, nil, fmt.Errorf("%s.source.query 不能为空", prefix)
			}
			if err := validateMapper(pipeline.SourceInterceptors, prefix); err != nil {
				return nil, nil, err
			}
			refs = append(refs, ref)
		case entity.SourceTypeAPI:
			ref := componentString(pipeline.Source, "connectionRef")
			if ref == "" {
				return nil, nil, fmt.Errorf("%s.source.connectionRef 不能为空", prefix)
			}
			if componentString(pipeline.Source, "path") == "" {
				return nil, nil, fmt.Errorf("%s.source.path 不能为空", prefix)
			}
			method := strings.ToUpper(componentString(pipeline.Source, "method"))
			if method != "GET" && method != "POST" {
				return nil, nil, fmt.Errorf("%s.source.method 只支持 GET 或 POST", prefix)
			}
			if pagination, ok := pipeline.Source["pagination"].(map[string]any); ok {
				t := componentString(pagination, "type")
				if t != "" && t != "none" && t != "page" && t != "offset" && t != "cursor" {
					return nil, nil, fmt.Errorf("%s.source.pagination.type 不支持 %q", prefix, t)
				}
			}
			if err := validateMapper(pipeline.SourceInterceptors, prefix); err != nil {
				return nil, nil, err
			}
			refs = append(refs, ref)
		default:
			return nil, nil, fmt.Errorf("不支持的 source_type %q", sourceType)
		}
	}
	return &bundle, refs, nil
}

func validateExcelParser(interceptors []dto.Component, prefix string) error {
	for _, interceptor := range interceptors {
		if componentString(interceptor, "type") == "parser" && componentString(interceptor, "parserType") == "document" {
			if value, ok := interceptor["xlsxMetricRule"]; ok && value != nil && fmt.Sprint(value) != "" {
				return nil
			}
		}
	}
	return fmt.Errorf("%s.sourceInterceptors 需要包含带 xlsxMetricRule 的 document parser", prefix)
}

func validateMapper(interceptors []dto.Component, prefix string) error {
	for _, interceptor := range interceptors {
		if componentString(interceptor, "type") != "mapper" {
			continue
		}
		mappings, ok := interceptor["mappings"].([]any)
		if !ok || len(mappings) == 0 {
			return fmt.Errorf("%s mapper.mappings 至少需要一项", prefix)
		}
		for i, item := range mappings {
			mapping, ok := item.(map[string]any)
			if !ok || componentString(mapping, "sourcePath") == "" || componentString(mapping, "targetField") == "" {
				return fmt.Errorf("%s mapper.mappings[%d] 需要 sourcePath 和 targetField", prefix, i)
			}
		}
		return nil
	}
	return fmt.Errorf("%s.sourceInterceptors 需要包含 mapper", prefix)
}

func componentString(component map[string]any, key string) string {
	value, ok := component[key]
	if !ok {
		return ""
	}
	text, ok := value.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(text)
}

func normalizeJSON(raw json.RawMessage) (string, error) {
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return "", err
	}
	normalized, err := json.Marshal(value)
	return string(normalized), err
}
