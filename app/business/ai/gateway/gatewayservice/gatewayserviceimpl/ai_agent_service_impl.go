package gatewayserviceimpl

import (
	"encoding/json"
	"errors"
	"strings"

	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"

	"github.com/gin-gonic/gin"
)

// AIAgentServiceImpl 提供智能体配置的业务实现。
type AIAgentServiceImpl struct {
	dao gatewaydao.IAIAgentDao
}

// NewAIAgentService 创建智能体配置服务。
func NewAIAgentService(dao gatewaydao.IAIAgentDao) gatewayservice.IAIAgentService {
	return &AIAgentServiceImpl{dao: dao}
}

// Create 新增智能体配置。
func (a *AIAgentServiceImpl) Create(c *gin.Context, req *gatewaymodels.AIAgentUpsert) (*gatewaymodels.AIAgent, error) {
	if err := a.prepareUpsert(req, false); err != nil {
		return nil, err
	}
	return a.dao.Create(c, req)
}

// Update 修改智能体配置。
func (a *AIAgentServiceImpl) Update(c *gin.Context, req *gatewaymodels.AIAgentUpsert) (*gatewaymodels.AIAgent, error) {
	if err := a.prepareUpsert(req, true); err != nil {
		return nil, err
	}
	current, err := a.dao.GetByID(c, req.ID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, errors.New("智能体不存在")
	}
	return a.dao.Update(c, req)
}

// DeleteByIDs 删除智能体配置。
func (a *AIAgentServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的智能体")
	}
	for _, id := range ids {
		if id == 0 {
			return errors.New("智能体ID不能为空")
		}
		current, err := a.dao.GetByID(c, id)
		if err != nil {
			return err
		}
		if current == nil {
			return errors.New("智能体不存在")
		}
	}
	return a.dao.DeleteByIDs(c, ids)
}

// GetByID 查询智能体详情。
func (a *AIAgentServiceImpl) GetByID(c *gin.Context, id int64) (*gatewaymodels.AIAgent, error) {
	if id == 0 {
		return nil, errors.New("id不能为空")
	}
	return a.dao.GetByID(c, id)
}

// List 查询智能体列表。
func (a *AIAgentServiceImpl) List(c *gin.Context, req *gatewaymodels.AIAgentQuery) (*gatewaymodels.AIAgentListData, error) {
	if req == nil {
		req = new(gatewaymodels.AIAgentQuery)
	}
	req.Name = strings.TrimSpace(req.Name)
	req.SandboxMode = strings.TrimSpace(req.SandboxMode)
	return a.dao.List(c, req)
}

func (a *AIAgentServiceImpl) prepareUpsert(req *gatewaymodels.AIAgentUpsert, isUpdate bool) error {
	if req == nil {
		return errors.New("参数不能为空")
	}
	if isUpdate && req.ID == 0 {
		return errors.New("id不能为空")
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Prompt = strings.TrimSpace(req.Prompt)
	req.DefaultLLMProviderID = strings.TrimSpace(req.DefaultLLMProviderID)
	req.DefaultLLMModelID = strings.TrimSpace(req.DefaultLLMModelID)
	req.SandboxMode = strings.TrimSpace(req.SandboxMode)
	req.WorkDir = strings.TrimSpace(req.WorkDir)
	req.MCPServerIDs = strings.TrimSpace(req.MCPServerIDs)
	req.MCPServerEnabledIDs = strings.TrimSpace(req.MCPServerEnabledIDs)
	if req.Name == "" {
		return errors.New("智能体名称不能为空")
	}
	if req.EnableLLMTemperature == nil {
		req.EnableLLMTemperature = agentBoolPtr(false)
	}
	if req.EnableLLMTopP == nil {
		req.EnableLLMTopP = agentBoolPtr(false)
	}
	if req.EnableLLMMaxTokens == nil {
		req.EnableLLMMaxTokens = agentBoolPtr(false)
	}
	if req.SandboxNetwork == nil {
		req.SandboxNetwork = agentBoolPtr(false)
	}
	if req.MCPEnabled == nil {
		req.MCPEnabled = agentBoolPtr(false)
	}
	if req.LLMMaxTokens < 0 || req.LLMMaxContextCount < 0 || req.RetrievalTopK < 0 {
		return errors.New("数字参数不能小于0")
	}
	if req.MCPServerIDs != "" {
		if err := validateAgentJSONArray(req.MCPServerIDs, "MCP服务ID列表"); err != nil {
			return err
		}
	}
	if req.MCPServerEnabledIDs != "" {
		if err := validateAgentJSONArray(req.MCPServerEnabledIDs, "启用MCP服务ID列表"); err != nil {
			return err
		}
	}
	return nil
}

func validateAgentJSONArray(content string, fieldName string) error {
	var data []interface{}
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return errors.New(fieldName + "必须是JSON数组字符串")
	}
	return nil
}

func agentBoolPtr(v bool) *bool {
	return &v
}
