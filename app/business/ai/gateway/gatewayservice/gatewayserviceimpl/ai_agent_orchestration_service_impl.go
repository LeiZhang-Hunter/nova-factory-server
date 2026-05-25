package gatewayserviceimpl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"nova-factory-server/app/constant/aiagent"
	"strconv"
	"strings"

	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"

	"github.com/gin-gonic/gin"
)

// AIAgentOrchestrationServiceImpl 提供智能体编排的业务实现。
type AIAgentOrchestrationServiceImpl struct {
	dao             gatewaydao.IAIAgentOrchestrationDao
	agentDao        gatewaydao.IAIAgentDao
	subAgentService gatewayservice.IAISubAgentService
}

// NewAIAgentOrchestrationService 创建智能体编排服务。
func NewAIAgentOrchestrationService(dao gatewaydao.IAIAgentOrchestrationDao,
	agentDao gatewaydao.IAIAgentDao, subAgentService gatewayservice.IAISubAgentService) gatewayservice.IAIAgentOrchestrationService {
	return &AIAgentOrchestrationServiceImpl{
		dao:             dao,
		agentDao:        agentDao,
		subAgentService: subAgentService,
	}
}

// Set 保存智能体编排配置。
func (a *AIAgentOrchestrationServiceImpl) Set(c *gin.Context, req *gatewaymodels.AIAgentOrchestrationUpsert) (*gatewaymodels.AIAgentOrchestration, error) {
	subAgents, err := a.prepareUpsert(c, req)
	if err != nil {
		return nil, err
	}

	agent, err := a.agentDao.GetByID(c, req.AgentID)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return nil, errors.New("智能体不存在")
	}

	current, err := a.dao.GetByAgentID(c, req.AgentID)
	if err != nil {
		return nil, err
	}

	var config gatewaymodels.AgentLoadConfig
	config.Agent = agent
	config.SubAgent = subAgents
	content, err := json.Marshal(&config)
	if err != nil {
		return nil, err
	}
	req.Config = string(content)
	if current != nil {
		return a.dao.UpdateByAgentID(c, req)
	}
	return a.dao.Create(c, req)
}

// Info 查询智能体编排详情。
func (a *AIAgentOrchestrationServiceImpl) Info(c *gin.Context, agentID int64) (*gatewaymodels.AIAgentOrchestration, error) {
	if agentID == 0 {
		return nil, errors.New("agentId不能为空")
	}
	if err := a.ensureAgentExists(c, agentID); err != nil {
		return nil, err
	}
	return a.dao.GetByAgentID(c, agentID)
}

// Remove 删除智能体编排配置。
func (a *AIAgentOrchestrationServiceImpl) Remove(c *gin.Context, agentIDs []int64) error {
	if len(agentIDs) == 0 {
		return errors.New("请选择要删除的编排配置")
	}
	for _, agentID := range agentIDs {
		if agentID == 0 {
			return errors.New("agentId不能为空")
		}
		if err := a.ensureAgentExists(c, agentID); err != nil {
			return err
		}
	}
	return a.dao.DeleteByAgentIDs(c, agentIDs)
}

func (a *AIAgentOrchestrationServiceImpl) prepareUpsert(c *gin.Context,
	req *gatewaymodels.AIAgentOrchestrationUpsert) ([]*gatewaymodels.AISubAgentUpsert, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	if req.AgentID == 0 {
		return nil, errors.New("agentId不能为空")
	}
	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		return nil, errors.New("content不能为空")
	}
	list, err := a.validateContent(c, req)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (a *AIAgentOrchestrationServiceImpl) ensureAgentExists(c *gin.Context, agentID int64) error {
	agent, err := a.agentDao.GetByID(c, agentID)
	if err != nil {
		return err
	}
	if agent == nil {
		return errors.New("智能体不存在")
	}
	return nil
}

// validateContent 组装配置
func (a *AIAgentOrchestrationServiceImpl) validateContent(c *gin.Context,
	req *gatewaymodels.AIAgentOrchestrationUpsert) ([]*gatewaymodels.AISubAgentUpsert, error) {
	var config gatewaymodels.AgentOrchestrationConfig
	if err := json.Unmarshal([]byte(req.Content), &config); err != nil {
		return nil, errors.New("content格式不正确")
	}

	if len(config.Edges) == 0 {
		return nil, errors.New("任务派发路线不能为空")
	}

	config.AgentId = strings.TrimSpace(config.AgentId)
	if config.AgentId == "" {
		return nil, errors.New("content.agentId不能为空")
	}
	contentAgentID, err := strconv.ParseInt(config.AgentId, 10, 64)
	if err != nil || contentAgentID == 0 {
		return nil, errors.New("content.agentId格式不正确")
	}
	if contentAgentID != req.AgentID {
		return nil, errors.New("content.agentId与agentId不一致")
	}
	if len(config.Nodes) == 0 {
		return nil, errors.New("content.nodes不能为空")
	}

	haveMaster := 0

	nodeMap := make(map[string]gatewaymodels.Node)
	dataMap := make(map[string]*gatewaymodels.AISubAgentUpsert)

	for index, node := range config.Nodes {
		nodeIndex := index + 1
		nodeMap[node.Id] = node
		if node.Type == aiagent.Master {
			haveMaster++
			continue
		}
		if strings.TrimSpace(node.Id) == "" {
			return nil, fmt.Errorf("content.nodes[%d].id不能为空", nodeIndex)
		}
		if strings.TrimSpace(node.Type) == "" {
			return nil, fmt.Errorf("content.nodes[%d].type不能为空", nodeIndex)
		}
		if node.Data.Config == nil {
			return nil, fmt.Errorf("content.nodes[%d].data.config不能为空", nodeIndex)
		}
		data, err := a.validateOrchestrationNodeConfig(c, node.Data.Config)
		if err != nil {
			return nil, err
		}

		if len(node.HandleBounds.Source) == 0 {
			return nil, errors.New(fmt.Sprintf("子Agent%s没有上游", node.Data.Config.Name))
		}

		dataMap[node.Id] = data
	}

	if haveMaster == 0 {
		return nil, errors.New("master 节点不存在")
	}

	if haveMaster > 1 {
		return nil, errors.New("master 只能存在一个")
	}

	if len(config.Nodes) <= 1 {
		return nil, errors.New("sub agent 节点没有")
	}

	list := make([]*gatewaymodels.AISubAgentUpsert, len(dataMap))

	//遍历所有边
	for i, edge := range config.Edges {
		_, ok := nodeMap[edge.Source]
		if !ok {
			return nil, errors.New(fmt.Sprintf(fmt.Sprintf("主Agent %s 不存在", edge.Source)))
		}

		_, ok = nodeMap[edge.Target]
		if !ok {
			return nil, errors.New(fmt.Sprintf("子Agent %s 不存在", edge.Target))
		}

		data, ok := dataMap[edge.Target]
		if !ok {
			continue
		}
		list[i] = data
	}

	return list, nil
}

func (a *AIAgentOrchestrationServiceImpl) validateOrchestrationNodeConfig(c *gin.Context, config *gatewaymodels.GraphAISubAgentUpsert) (*gatewaymodels.AISubAgentUpsert, error) {
	data := &gatewaymodels.AISubAgentUpsert{
		ID:                        config.ID,
		Name:                      config.Name,
		Type:                      config.SubAgentType,
		CoreSubAgent:              config.CoreSubAgent,
		Description:               config.Description,
		Instruction:               config.Instruction,
		MCPEnabled:                config.MCPEnabled,
		MCPServerIDs:              config.MCPServerIDs,
		MCPServerEnabledIDs:       config.MCPServerEnabledIDs,
		LocalToolEnabled:          config.LocalToolEnabled,
		LocalTools:                config.LocalTools,
		AllowMcpServerIdsToolsRaw: config.AllowMcpServerIdsToolsRaw,
		AllowMcpServerIdsTools:    config.AllowMcpServerIdsTools,
		Enable:                    config.Enable,
	}
	err := a.subAgentService.ValidateType(c, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetConfigInfo 查询智能体编排详情。
func (a *AIAgentOrchestrationServiceImpl) GetConfigInfo(c context.Context, agentID int64) (*gatewaymodels.AIAgentOrchestration, error) {
	return a.dao.GetConfigByAgentID(c, agentID)
}
