package impl

import (
	"context"
	"encoding/json"
	"testing"

	"nova-factory-server/app/business/data/models/dto"
	"nova-factory-server/app/business/data/models/entity"

	"github.com/gin-gonic/gin"
)

func raw(value any) json.RawMessage {
	payload, _ := json.Marshal(value)
	return payload
}

func bundle(sourceType string, source map[string]any, interceptor map[string]any) json.RawMessage {
	return raw(map[string]any{
		"apiVersion": "nova.data/v1",
		"kind":       "PipelineBundle",
		"pipelines": []any{map[string]any{
			"name":               sourceType + "-pipeline",
			"source":             source,
			"queue":              map[string]any{"type": "channel", "name": "main", "capacity": 10},
			"ontologySinks":      []any{map[string]any{"type": "memory", "name": "memory"}},
			"sourceInterceptors": []any{interceptor},
			"sinkInterceptors":   []any{},
		}},
	})
}

func TestValidatePipelineBundleSupportedTypes(t *testing.T) {
	tests := []struct {
		name       string
		sourceType string
		config     json.RawMessage
		wantRefs   int
	}{
		{"file", entity.SourceTypeFile, bundle("file", map[string]any{"type": "file", "name": "files", "path": "/data/*.xlsx"}, map[string]any{"type": "parser", "name": "parser", "parserType": "document", "xlsxMetricRule": map[string]any{"matchRule": map[string]any{}}}), 0},
		{"mysql", entity.SourceTypeMySQL, bundle("mysql", map[string]any{"type": "mysql", "name": "db", "connectionRef": "mysql-1", "query": "select 1"}, map[string]any{"type": "mapper", "name": "mapper", "mappings": []any{map[string]any{"sourcePath": "id", "targetField": "id"}}}), 1},
		{"api", entity.SourceTypeAPI, bundle("api", map[string]any{"type": "api", "name": "api", "connectionRef": "api-1", "path": "/orders", "method": "GET", "pagination": map[string]any{"type": "page"}}, map[string]any{"type": "mapper", "name": "mapper", "mappings": []any{map[string]any{"sourcePath": "$.id", "targetField": "id"}}}), 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, refs, err := validatePipelineBundle(tt.config, tt.sourceType)
			if err != nil {
				t.Fatalf("validate: %v", err)
			}
			if len(refs) != tt.wantRefs {
				t.Fatalf("refs=%d, want %d", len(refs), tt.wantRefs)
			}
		})
	}
}

func TestValidatePipelineBundleRejectsMixedAndDuplicatePipelines(t *testing.T) {
	config := bundle("mysql", map[string]any{"type": "mysql", "name": "db", "connectionRef": "mysql-1", "query": "select 1"}, map[string]any{"type": "mapper", "name": "mapper", "mappings": []any{map[string]any{"sourcePath": "id", "targetField": "id"}}})
	if _, _, err := validatePipelineBundle(config, entity.SourceTypeAPI); err == nil {
		t.Fatal("expected mixed source type error")
	}

	var value map[string]any
	_ = json.Unmarshal(config, &value)
	pipelines := value["pipelines"].([]any)
	value["pipelines"] = append(pipelines, pipelines[0])
	if _, _, err := validatePipelineBundle(raw(value), entity.SourceTypeMySQL); err == nil {
		t.Fatal("expected duplicate pipeline name error")
	}
}

type ruleDAOStub struct {
	rule       *entity.PipelineRule
	references int64
}

func (s *ruleDAOStub) Create(context.Context, *entity.PipelineRule) error { return nil }
func (s *ruleDAOStub) GetByID(context.Context, string) (*entity.PipelineRule, error) {
	return s.rule, nil
}
func (s *ruleDAOStub) GetByName(context.Context, string, string) (*entity.PipelineRule, error) {
	return nil, nil
}
func (s *ruleDAOStub) List(context.Context, int, int, string, string, string) ([]entity.PipelineRule, int64, error) {
	return nil, 0, nil
}
func (s *ruleDAOStub) Update(_ context.Context, rule *entity.PipelineRule) error {
	s.rule = rule
	return nil
}
func (s *ruleDAOStub) Delete(context.Context, string) error { return nil }
func (s *ruleDAOStub) CountConnectionReferences(context.Context, string) (int64, error) {
	return s.references, nil
}

type connectionDAOStub struct {
	rows map[string]*entity.ServiceConnection
}

func (s *connectionDAOStub) Create(_ context.Context, row *entity.ServiceConnection) error {
	if s.rows == nil {
		s.rows = map[string]*entity.ServiceConnection{}
	}
	s.rows[row.ID] = row
	return nil
}
func (s *connectionDAOStub) GetByID(_ context.Context, id string) (*entity.ServiceConnection, error) {
	return s.rows[id], nil
}
func (s *connectionDAOStub) GetByName(context.Context, string, string) (*entity.ServiceConnection, error) {
	return nil, nil
}
func (s *connectionDAOStub) List(context.Context, int, int, string, string, string) ([]entity.ServiceConnection, int64, error) {
	return nil, 0, nil
}
func (s *connectionDAOStub) Update(_ context.Context, row *entity.ServiceConnection) error {
	s.rows[row.ID] = row
	return nil
}
func (s *connectionDAOStub) Delete(_ context.Context, id string) error {
	delete(s.rows, id)
	return nil
}

func TestPipelineRuleUpdateValidatesReferenceAndIncrementsVersion(t *testing.T) {
	config := bundle("mysql", map[string]any{"type": "mysql", "name": "db", "connectionRef": "mysql-1", "query": "select 1"}, map[string]any{"type": "mapper", "name": "mapper", "mappings": []any{map[string]any{"sourcePath": "id", "targetField": "id"}}})
	rules := &ruleDAOStub{rule: &entity.PipelineRule{ID: "rule-1", Name: "old", SourceType: "mysql", Status: "enabled", Version: 2, Config: string(config)}}
	connections := &connectionDAOStub{rows: map[string]*entity.ServiceConnection{"mysql-1": {ID: "mysql-1", SourceType: "mysql"}}}
	service := &PipelineRuleService{rules: rules, connections: connections}
	result, err := service.Update(&gin.Context{}, "rule-1", &dto.UpdatePipelineRuleRequest{Name: "new", Status: "enabled", Config: config})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if result.Version != 3 {
		t.Fatalf("version=%d, want 3", result.Version)
	}

	delete(connections.rows, "mysql-1")
	if _, err := service.Update(&gin.Context{}, "rule-1", &dto.UpdatePipelineRuleRequest{Name: "new", Status: "enabled", Config: config}); err == nil {
		t.Fatal("expected invalid reference error")
	}
}

func TestConnectionCredentialsEncryptedAndMasked(t *testing.T) {
	connections := &connectionDAOStub{rows: map[string]*entity.ServiceConnection{}}
	service := &ServiceConnectionService{connections: connections, rules: &ruleDAOStub{}, key: []byte("0123456789abcdef")}
	result, err := service.Create(&gin.Context{}, &dto.CreateServiceConnectionRequest{Name: "erp", SourceType: "api", Status: "enabled", Config: raw(map[string]any{"baseURL": "https://example.com", "authType": "bearer"}), Credentials: map[string]string{"token": "plain-secret"}})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	stored := connections.rows[result.ID]
	if stored.EncryptedCredentials == "" || stored.EncryptedCredentials == "plain-secret" {
		t.Fatal("credentials were not encrypted")
	}
	if !result.CredentialsConfigured || len(result.CredentialFields) != 1 || result.CredentialFields[0] != "token" {
		t.Fatalf("unexpected masked response: %+v", result)
	}
}
