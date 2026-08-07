package impl

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"nova-factory-server/app/business/data/dao"
	"nova-factory-server/app/business/data/models/dto"
	"nova-factory-server/app/business/data/models/entity"
	"nova-factory-server/app/business/data/service"
	"nova-factory-server/app/utils/aes"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type PipelineRuleService struct {
	rules       dao.IPipelineRuleDAO
	connections dao.IServiceConnectionDAO
}

type ServiceConnectionService struct {
	connections dao.IServiceConnectionDAO
	rules       dao.IPipelineRuleDAO
	key         []byte
}

func NewPipelineRuleService(rules dao.IPipelineRuleDAO, connections dao.IServiceConnectionDAO) service.IPipelineRuleService {
	return &PipelineRuleService{rules: rules, connections: connections}
}

func NewServiceConnectionService(connections dao.IServiceConnectionDAO, rules dao.IPipelineRuleDAO) service.IServiceConnectionService {
	keyText := strings.TrimSpace(viper.GetString("data.credentials_key"))
	if keyText == "" {
		keyText = strings.TrimSpace(viper.GetString("aes_key"))
	}
	key, _ := aes.DecodeKeyString(keyText)
	return &ServiceConnectionService{connections: connections, rules: rules, key: key}
}

func (s *PipelineRuleService) Create(c *gin.Context, req *dto.CreatePipelineRuleRequest) (*dto.PipelineRuleResponse, error) {
	if err := s.ensureUnique(c, req.SourceType, req.Name, ""); err != nil {
		return nil, err
	}
	config, refs, err := validateAndNormalize(req.Config, req.SourceType)
	if err != nil {
		return nil, err
	}
	if err := s.validateReferences(c, req.SourceType, refs); err != nil {
		return nil, err
	}
	status := req.Status
	if status == "" {
		status = entity.RuleStatusEnabled
	}
	rule := &entity.PipelineRule{ID: uuid.NewString(), Name: req.Name, Description: req.Description, SourceType: req.SourceType, Status: status, Version: 1, Config: config}
	if err := s.rules.Create(c, rule); err != nil {
		return nil, err
	}
	return pipelineRuleResponse(rule), nil
}

func (s *PipelineRuleService) Get(c *gin.Context, id string) (*dto.PipelineRuleResponse, error) {
	rule, err := s.rules.GetByID(c, id)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, fmt.Errorf("Pipeline 规则不存在")
	}
	return pipelineRuleResponse(rule), nil
}

func (s *PipelineRuleService) List(c *gin.Context, offset, limit int, sourceType, name, status string) ([]dto.PipelineRuleResponse, int64, error) {
	rows, total, err := s.rules.List(c, offset, limit, sourceType, name, status)
	if err != nil {
		return nil, 0, err
	}
	out := make([]dto.PipelineRuleResponse, 0, len(rows))
	for i := range rows {
		out = append(out, *pipelineRuleResponse(&rows[i]))
	}
	return out, total, nil
}

func (s *PipelineRuleService) Update(c *gin.Context, id string, req *dto.UpdatePipelineRuleRequest) (*dto.PipelineRuleResponse, error) {
	rule, err := s.rules.GetByID(c, id)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, fmt.Errorf("Pipeline 规则不存在")
	}
	if err := s.ensureUnique(c, rule.SourceType, req.Name, id); err != nil {
		return nil, err
	}
	config, refs, err := validateAndNormalize(req.Config, rule.SourceType)
	if err != nil {
		return nil, err
	}
	if err := s.validateReferences(c, rule.SourceType, refs); err != nil {
		return nil, err
	}
	rule.Name, rule.Description, rule.Status, rule.Config = req.Name, req.Description, req.Status, config
	rule.Version++
	if err := s.rules.Update(c, rule); err != nil {
		return nil, err
	}
	return pipelineRuleResponse(rule), nil
}

func (s *PipelineRuleService) Delete(c *gin.Context, id string) error {
	rule, err := s.rules.GetByID(c, id)
	if err != nil {
		return err
	}
	if rule == nil {
		return fmt.Errorf("Pipeline 规则不存在")
	}
	return s.rules.Delete(c, id)
}

func (s *PipelineRuleService) ValidateConfig(c *gin.Context, raw json.RawMessage, sourceType string) (json.RawMessage, error) {
	config, refs, err := validateAndNormalize(raw, sourceType)
	if err != nil {
		return nil, err
	}
	if err := s.validateReferences(c, sourceType, refs); err != nil {
		return nil, err
	}
	return json.RawMessage(config), nil
}

func (s *PipelineRuleService) ensureUnique(c *gin.Context, sourceType, name, exceptID string) error {
	existing, err := s.rules.GetByName(c, sourceType, name)
	if err != nil {
		return err
	}
	if existing != nil && existing.ID != exceptID {
		return fmt.Errorf("%s 类型下已存在名为 %q 的规则", sourceType, name)
	}
	return nil
}

func (s *PipelineRuleService) validateReferences(c *gin.Context, sourceType string, refs []string) error {
	seen := make(map[string]struct{})
	for _, ref := range refs {
		if _, ok := seen[ref]; ok {
			continue
		}
		seen[ref] = struct{}{}
		connection, err := s.connections.GetByID(c, ref)
		if err != nil {
			return err
		}
		if connection == nil {
			return fmt.Errorf("连接档案 %q 不存在", ref)
		}
		if connection.SourceType != sourceType {
			return fmt.Errorf("连接档案 %q 的类型不是 %s", ref, sourceType)
		}
	}
	return nil
}

func validateAndNormalize(raw json.RawMessage, sourceType string) (string, []string, error) {
	_, refs, err := validatePipelineBundle(raw, sourceType)
	if err != nil {
		return "", nil, err
	}
	config, err := normalizeJSON(raw)
	return config, refs, err
}

func pipelineRuleResponse(rule *entity.PipelineRule) *dto.PipelineRuleResponse {
	return &dto.PipelineRuleResponse{ID: rule.ID, Name: rule.Name, Description: rule.Description, SourceType: rule.SourceType, Status: rule.Status, Version: rule.Version, Config: json.RawMessage(rule.Config), CreatedAt: rule.CreatedAt, UpdatedAt: rule.UpdatedAt}
}

func (s *ServiceConnectionService) Create(c *gin.Context, req *dto.CreateServiceConnectionRequest) (*dto.ServiceConnectionResponse, error) {
	if err := s.ensureUnique(c, req.SourceType, req.Name, ""); err != nil {
		return nil, err
	}
	config, err := validateConnectionConfig(req.SourceType, req.Config)
	if err != nil {
		return nil, err
	}
	encrypted, err := s.encryptCredentials(req.Credentials)
	if err != nil {
		return nil, err
	}
	status := req.Status
	if status == "" {
		status = entity.RuleStatusEnabled
	}
	row := &entity.ServiceConnection{ID: uuid.NewString(), Name: req.Name, Description: req.Description, SourceType: req.SourceType, Status: status, Config: config, EncryptedCredentials: encrypted}
	if err := s.connections.Create(c, row); err != nil {
		return nil, err
	}
	return s.connectionResponse(row), nil
}

func (s *ServiceConnectionService) Get(c *gin.Context, id string) (*dto.ServiceConnectionResponse, error) {
	row, err := s.connections.GetByID(c, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, fmt.Errorf("连接档案不存在")
	}
	return s.connectionResponse(row), nil
}

func (s *ServiceConnectionService) List(c *gin.Context, offset, limit int, sourceType, name, status string) ([]dto.ServiceConnectionResponse, int64, error) {
	rows, total, err := s.connections.List(c, offset, limit, sourceType, name, status)
	if err != nil {
		return nil, 0, err
	}
	out := make([]dto.ServiceConnectionResponse, 0, len(rows))
	for i := range rows {
		out = append(out, *s.connectionResponse(&rows[i]))
	}
	return out, total, nil
}

func (s *ServiceConnectionService) Update(c *gin.Context, id string, req *dto.UpdateServiceConnectionRequest) (*dto.ServiceConnectionResponse, error) {
	row, err := s.connections.GetByID(c, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, fmt.Errorf("连接档案不存在")
	}
	if err := s.ensureUnique(c, row.SourceType, req.Name, id); err != nil {
		return nil, err
	}
	config, err := validateConnectionConfig(row.SourceType, req.Config)
	if err != nil {
		return nil, err
	}
	row.Name, row.Description, row.Status, row.Config = req.Name, req.Description, req.Status, config
	if req.Credentials != nil {
		row.EncryptedCredentials, err = s.encryptCredentials(*req.Credentials)
		if err != nil {
			return nil, err
		}
	}
	if err := s.connections.Update(c, row); err != nil {
		return nil, err
	}
	return s.connectionResponse(row), nil
}

func (s *ServiceConnectionService) Delete(c *gin.Context, id string) error {
	row, err := s.connections.GetByID(c, id)
	if err != nil {
		return err
	}
	if row == nil {
		return fmt.Errorf("连接档案不存在")
	}
	count, err := s.rules.CountConnectionReferences(c, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("连接档案已被 %d 条 Pipeline 规则引用，不能删除", count)
	}
	return s.connections.Delete(c, id)
}

func (s *ServiceConnectionService) ensureUnique(c *gin.Context, sourceType, name, exceptID string) error {
	existing, err := s.connections.GetByName(c, sourceType, name)
	if err != nil {
		return err
	}
	if existing != nil && existing.ID != exceptID {
		return fmt.Errorf("%s 类型下已存在名为 %q 的连接档案", sourceType, name)
	}
	return nil
}

func validateConnectionConfig(sourceType string, raw json.RawMessage) (string, error) {
	var config map[string]any
	if err := json.Unmarshal(raw, &config); err != nil {
		return "", fmt.Errorf("连接配置不是有效 JSON: %w", err)
	}
	switch sourceType {
	case entity.SourceTypeMySQL:
		if componentString(config, "host") == "" || componentString(config, "database") == "" || componentString(config, "username") == "" {
			return "", fmt.Errorf("MySQL 连接需要 host、database 和 username")
		}
	case entity.SourceTypeAPI:
		if componentString(config, "baseURL") == "" {
			return "", fmt.Errorf("API 连接需要 baseURL")
		}
		authType := componentString(config, "authType")
		if authType != "" && authType != "none" && authType != "basic" && authType != "bearer" && authType != "apiKey" {
			return "", fmt.Errorf("不支持的 API authType %q", authType)
		}
	default:
		return "", fmt.Errorf("连接档案只支持 mysql 或 api")
	}
	return normalizeJSON(raw)
}

func (s *ServiceConnectionService) encryptCredentials(credentials map[string]string) (string, error) {
	if len(credentials) == 0 {
		return "", nil
	}
	if len(s.key) == 0 {
		return "", fmt.Errorf("未配置 data.credentials_key")
	}
	payload, err := json.Marshal(credentials)
	if err != nil {
		return "", err
	}
	return aes.EncryptString(s.key, string(payload))
}

func (s *ServiceConnectionService) credentialFields(encrypted string) []string {
	if encrypted == "" || len(s.key) == 0 {
		return []string{}
	}
	plain, err := aes.DecryptString(s.key, encrypted)
	if err != nil {
		return []string{}
	}
	var credentials map[string]string
	if json.Unmarshal([]byte(plain), &credentials) != nil {
		return []string{}
	}
	fields := make([]string, 0, len(credentials))
	for key, value := range credentials {
		if value != "" {
			fields = append(fields, key)
		}
	}
	sort.Strings(fields)
	return fields
}

func (s *ServiceConnectionService) connectionResponse(row *entity.ServiceConnection) *dto.ServiceConnectionResponse {
	fields := s.credentialFields(row.EncryptedCredentials)
	return &dto.ServiceConnectionResponse{ID: row.ID, Name: row.Name, Description: row.Description, SourceType: row.SourceType, Status: row.Status, Config: json.RawMessage(row.Config), CredentialFields: fields, CredentialsConfigured: len(fields) > 0, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt}
}
