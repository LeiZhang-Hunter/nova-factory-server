package impl

import (
	"encoding/json"
	"fmt"
	"time"

	"nova-factory-server/app/business/data/dao"
	"nova-factory-server/app/business/data/models/dto"
	"nova-factory-server/app/business/data/models/entity"
	"nova-factory-server/app/business/data/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ICollectorRuleServiceImpl 采集器规则下发服务实现
type ICollectorRuleServiceImpl struct {
	db           *gorm.DB
	collectors   dao.ICollectorDAO
	rules        dao.IPipelineRuleDAO
	dispatchLogs dao.IDispatchLogDAO
}

// NewICollectorRuleServiceImpl 创建采集器规则下发服务
func NewICollectorRuleServiceImpl(
	db *gorm.DB,
	collectors dao.ICollectorDAO,
	rules dao.IPipelineRuleDAO,
	dispatchLogs dao.IDispatchLogDAO,
) service.ICollectorRuleService {
	return &ICollectorRuleServiceImpl{
		db:           db,
		collectors:   collectors,
		rules:        rules,
		dispatchLogs: dispatchLogs,
	}
}

// Dispatch 一次下发：多选规则合并为一条配置快照，写入一条下发记录
func (s *ICollectorRuleServiceImpl) Dispatch(c *gin.Context, collectorID, operator string, req *dto.BindCollectorRulesRequest) (*dto.DispatchLogResponse, error) {
	col, err := s.collectors.GetByID(c, collectorID)
	if err != nil {
		return nil, err
	}
	if col == nil {
		return nil, fmt.Errorf("采集器不存在")
	}
	rules, err := s.validateRules(c, req.RuleIDs)
	if err != nil {
		return nil, err
	}
	payload, checksum, err := buildMergedConfig(rules)
	if err != nil {
		return nil, err
	}
	ruleIDs := make([]string, 0, len(rules))
	ruleNames := make([]string, 0, len(rules))
	for i := range rules {
		ruleIDs = append(ruleIDs, rules[i].ID)
		ruleNames = append(ruleNames, rules[i].Name)
	}
	idsJSON, err := json.Marshal(ruleIDs)
	if err != nil {
		return nil, err
	}
	namesJSON, err := json.Marshal(ruleNames)
	if err != nil {
		return nil, err
	}
	no, err := s.dispatchLogs.NextDispatchNo(c, collectorID)
	if err != nil {
		return nil, err
	}
	log := &entity.DispatchLog{
		ID:            uuid.NewString(),
		CollectorID:   collectorID,
		CollectorName: col.Name,
		RuleIDs:       string(idsJSON),
		RuleNames:     string(namesJSON),
		DispatchNo:    no,
		Config:        string(payload),
		ConfigMD5:     checksum,
		CreatedBy:     operator,
	}
	err = s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		return s.dispatchLogs.CreateWithTx(c, tx, log)
	})
	if err != nil {
		return nil, err
	}
	return dispatchLogResponse(log)
}

// PreviewProposed 拟下发集合预览（dry-run，不落库）
func (s *ICollectorRuleServiceImpl) PreviewProposed(c *gin.Context, collectorID string, req *dto.BindCollectorRulesRequest) (*dto.DispatchPreviewResponse, error) {
	col, err := s.collectors.GetByID(c, collectorID)
	if err != nil {
		return nil, err
	}
	if col == nil {
		return nil, fmt.Errorf("采集器不存在")
	}
	rules, err := s.validateRules(c, req.RuleIDs)
	if err != nil {
		return nil, err
	}
	payload, checksum, err := buildMergedConfig(rules)
	if err != nil {
		return nil, err
	}
	return previewResponse(col.ID, rules, string(payload), checksum), nil
}

// PreviewCurrent 最近一次下发记录的配置预览
func (s *ICollectorRuleServiceImpl) PreviewCurrent(c *gin.Context, collectorID string) (*dto.DispatchPreviewResponse, error) {
	col, err := s.collectors.GetByID(c, collectorID)
	if err != nil {
		return nil, err
	}
	if col == nil {
		return nil, fmt.Errorf("采集器不存在")
	}
	latest, err := s.dispatchLogs.LatestByCollector(c, collectorID)
	if err != nil {
		return nil, err
	}
	if latest == nil {
		payload, checksum, err := buildMergedConfig(nil)
		if err != nil {
			return nil, err
		}
		return previewResponse(col.ID, nil, string(payload), checksum), nil
	}
	rules, err := s.rulesByIDs(c, unmarshalNames(latest.RuleIDs))
	if err != nil {
		return nil, err
	}
	return previewResponse(col.ID, rules, latest.Config, latest.ConfigMD5), nil
}

// PullRules 采集器拉取最近一次下发的合并配置，并顺带刷新在线状态
func (s *ICollectorRuleServiceImpl) PullRules(c *gin.Context, deviceID, token, clientMD5 string) (*dto.PullRulesResponse, error) {
	col, err := s.collectors.GetByDeviceID(c, deviceID)
	if err != nil {
		return nil, err
	}
	if col == nil {
		return nil, fmt.Errorf("设备认证失败")
	}
	latest, err := s.dispatchLogs.LatestByCollector(c, col.ID)
	if err != nil {
		return nil, err
	}
	config := ""
	checksum := ""
	if latest != nil {
		config = latest.Config
		checksum = latest.ConfigMD5
	} else {
		_, emptyChecksum, err := buildMergedConfig(nil)
		if err != nil {
			return nil, err
		}
		checksum = emptyChecksum
	}
	updated := clientMD5 != checksum
	now := time.Now()
	var dispatchedAt *time.Time
	lastChecksum := ""
	if updated {
		dispatchedAt = &now
		lastChecksum = checksum
	}
	if err := s.collectors.UpdateHeartbeat(c, col.ID, now, now, lastChecksum, dispatchedAt); err != nil {
		return nil, err
	}
	resp := &dto.PullRulesResponse{DeviceID: deviceID, Checksum: checksum, Updated: updated}
	if updated {
		resp.Config = config
	}
	return resp, nil
}

// ListLogs 分页查询下发记录
func (s *ICollectorRuleServiceImpl) ListLogs(c *gin.Context, offset, limit int, collectorID string) ([]dto.DispatchLogResponse, int64, error) {
	rows, total, err := s.dispatchLogs.List(c, offset, limit, collectorID)
	if err != nil {
		return nil, 0, err
	}
	out := make([]dto.DispatchLogResponse, 0, len(rows))
	for i := range rows {
		item, err := dispatchLogResponse(&rows[i])
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *item)
	}
	return out, total, nil
}

// validateRules 校验规则存在且启用
func (s *ICollectorRuleServiceImpl) validateRules(c *gin.Context, ruleIDs []string) ([]entity.PipelineRule, error) {
	rules, err := s.rulesByIDs(c, ruleIDs)
	if err != nil {
		return nil, err
	}
	found := make(map[string]*entity.PipelineRule, len(rules))
	for i := range rules {
		found[rules[i].ID] = &rules[i]
	}
	out := make([]entity.PipelineRule, 0, len(ruleIDs))
	for _, id := range ruleIDs {
		rule, ok := found[id]
		if !ok {
			return nil, fmt.Errorf("规则 %s 不存在", id)
		}
		if rule.Status != entity.RuleStatusEnabled {
			return nil, fmt.Errorf("规则 %s 已禁用，无法下发", rule.Name)
		}
		out = append(out, *rule)
	}
	return out, nil
}

// rulesByIDs 按 ID 查询规则
func (s *ICollectorRuleServiceImpl) rulesByIDs(c *gin.Context, ids []string) ([]entity.PipelineRule, error) {
	return s.rules.ListByIDs(c, ids)
}

func unmarshalNames(raw string) []string {
	if raw == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}

func dispatchLogResponse(log *entity.DispatchLog) (*dto.DispatchLogResponse, error) {
	return &dto.DispatchLogResponse{
		ID:            log.ID,
		CollectorID:   log.CollectorID,
		CollectorName: log.CollectorName,
		RuleIDs:       unmarshalNames(log.RuleIDs),
		RuleNames:     unmarshalNames(log.RuleNames),
		DispatchNo:    log.DispatchNo,
		Config:        log.Config,
		ConfigMD5:     log.ConfigMD5,
		CreatedAt:     log.CreatedAt,
		CreatedBy:     log.CreatedBy,
	}, nil
}

func previewResponse(collectorID string, rules []entity.PipelineRule, config string, checksum string) *dto.DispatchPreviewResponse {
	items := make([]dto.PreviewRuleItem, 0, len(rules))
	for i := range rules {
		items = append(items, dto.PreviewRuleItem{
			RuleID:     rules[i].ID,
			RuleName:   rules[i].Name,
			SourceType: rules[i].SourceType,
			Status:     rules[i].Status,
			Version:    rules[i].Version,
		})
	}
	return &dto.DispatchPreviewResponse{
		CollectorID: collectorID,
		Rules:       items,
		Checksum:    checksum,
		Config:      config,
	}
}
