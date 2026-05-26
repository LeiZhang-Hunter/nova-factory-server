package gatewayserviceimpl

import (
	"encoding/json"
	"errors"
	"fmt"
	"nova-factory-server/app/business/admin/system/systemdao"
	"nova-factory-server/app/constant/aiagent"
	"strconv"
	"strings"

	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"

	"github.com/gin-gonic/gin"
)

// AISubAgentServiceImpl 提供子智能体配置的业务实现。
type AISubAgentServiceImpl struct {
	dao          gatewaydao.IAISubAgentDao
	mcpServerDao gatewaydao.IMCPServerDao
	dictDataDao  systemdao.IDictDataDao
}

// NewAISubAgentService 创建子智能体配置服务。
func NewAISubAgentService(dao gatewaydao.IAISubAgentDao,
	mcpServerDao gatewaydao.IMCPServerDao,
	dictDataDao systemdao.IDictDataDao) gatewayservice.IAISubAgentService {
	return &AISubAgentServiceImpl{
		dao:          dao,
		mcpServerDao: mcpServerDao,
		dictDataDao:  dictDataDao,
	}
}

// ValidateType 子agent类型
func (a *AISubAgentServiceImpl) ValidateType(c *gin.Context, req *gatewaymodels.AISubAgentUpsert) error {
	if req == nil {
		return errors.New("参数不能为空")
	}
	req.Type = strings.TrimSpace(req.Type)
	if req.Type == "" {
		return errors.New("type不能为空")
	}
	typeKey := strings.ToLower(req.Type)
	rows := a.dictDataDao.SelectDictDataByType(c, aiagent.SubAgentType)
	var checkDict bool
	for _, row := range rows {
		if row == nil {
			continue
		}
		if strings.ToLower(strings.TrimSpace(row.DictValue)) == typeKey {
			checkDict = true
		}
	}

	if !checkDict {
		return errors.New("agent类型设置失败")
	}

	checkDict = false

	if req.Type == aiagent.CORE {
		if req.CoreSubAgent == "" {
			return errors.New("核心子agent不能为空")
		}

		rows := a.dictDataDao.SelectDictDataByType(c, aiagent.CoreSubAgent)
		for _, row := range rows {
			if row == nil {
				continue
			}
			if strings.ToLower(strings.TrimSpace(row.DictValue)) == typeKey {
				checkDict = true
			}
		}

		if !checkDict {
			return errors.New("核心agent不存在")
		}
		return nil
	}

	// 自定义Agent校验字段
	if req.Name == "" {
		return errors.New("agent名字不是空的")
	}

	if req.Description == "" {
		return errors.New("agent描述不能为空")
	}

	if req.Instruction == "" {
		return errors.New("agent指令不能为空")
	}
	return nil
}

// Create 新增子智能体配置。
func (a *AISubAgentServiceImpl) Create(c *gin.Context, req *gatewaymodels.AISubAgentUpsert) (*gatewaymodels.AISubAgent, error) {
	if err := a.ValidateType(c, req); err != nil {
		return nil, err
	}
	//检查名字是否存在
	info, err := a.dao.GetByName(c, req.Name)
	if err != nil {
		return nil, err
	}
	if info != nil {
		return nil, errors.New("智能体名称已经存在，请切换一个智能体")
	}
	if err := a.prepareUpsert(c, req, false); err != nil {
		return nil, err
	}
	return a.dao.Create(c, req)
}

// Update 修改子智能体配置。
func (a *AISubAgentServiceImpl) Update(c *gin.Context, req *gatewaymodels.AISubAgentUpsert) (*gatewaymodels.AISubAgent, error) {
	if err := a.ValidateType(c, req); err != nil {
		return nil, err
	}
	if err := a.prepareUpsert(c, req, true); err != nil {
		return nil, err
	}
	current, err := a.dao.GetByID(c, req.ID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, errors.New("子智能体不存在")
	}
	return a.dao.Update(c, req)
}

// DeleteByIDs 删除子智能体配置。
func (a *AISubAgentServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的子智能体")
	}
	for _, id := range ids {
		if id == 0 {
			return errors.New("子智能体ID不能为空")
		}
		current, err := a.dao.GetByID(c, id)
		if err != nil {
			return err
		}
		if current == nil {
			return errors.New("子智能体不存在")
		}
	}
	return a.dao.DeleteByIDs(c, ids)
}

// GetByID 查询子智能体详情。
func (a *AISubAgentServiceImpl) GetByID(c *gin.Context, id int64) (*gatewaymodels.AISubAgent, error) {
	if id == 0 {
		return nil, errors.New("id不能为空")
	}
	return a.dao.GetByID(c, id)
}

// List 查询子智能体列表。
func (a *AISubAgentServiceImpl) List(c *gin.Context, req *gatewaymodels.AISubAgentQuery) (*gatewaymodels.AISubAgentListData, error) {
	if req == nil {
		req = new(gatewaymodels.AISubAgentQuery)
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Type = strings.TrimSpace(req.Type)
	return a.dao.List(c, req)
}

func (a *AISubAgentServiceImpl) prepareUpsert(c *gin.Context, req *gatewaymodels.AISubAgentUpsert, isUpdate bool) error {
	if req == nil {
		return errors.New("参数不能为空")
	}
	if isUpdate && req.ID == 0 {
		return errors.New("id不能为空")
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Type = strings.TrimSpace(req.Type)
	req.Description = strings.TrimSpace(req.Description)
	req.Instruction = strings.TrimSpace(req.Instruction)
	if req.Name == "" {
		return errors.New("子智能体名称不能为空")
	}
	if req.Type == "" {
		return errors.New("子智能体类型不能为空")
	}
	if req.MCPEnabled == nil {
		req.MCPEnabled = boolPtr(false)
	}
	if req.LocalToolEnabled == nil {
		req.LocalToolEnabled = boolPtr(false)
	}
	if req.Enable == nil {
		req.Enable = boolPtr(false)
	}

	return a.prepareAllowMcpServerIdsTools(c, req)
}

func (a *AISubAgentServiceImpl) prepareAllowMcpServerIdsTools(c *gin.Context, req *gatewaymodels.AISubAgentUpsert) error {
	if req == nil {
		return errors.New("参数不能为空")
	}
	if len(req.AllowMcpServerIdsTools) == 0 {
		req.AllowMcpServerIdsToolsRaw = ""
		return nil
	}

	normalized := make(map[string][]string, len(req.AllowMcpServerIdsTools))
	for serverIDText, tools := range req.AllowMcpServerIdsTools {
		serverIDText = strings.TrimSpace(serverIDText)
		if serverIDText == "" {
			return errors.New("允许使用工具的MCP服务ID不能为空")
		}
		serverID, err := strconv.ParseInt(serverIDText, 10, 64)
		if err != nil || serverID <= 0 {
			return fmt.Errorf("MCP服务ID格式不正确: %s", serverIDText)
		}
		server, err := a.mcpServerDao.GetByID(c, serverID)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("MCP服务不存在: %s", serverIDText)
		}

		toolNames := make([]string, 0, len(tools))
		for _, tool := range tools {
			tool = strings.TrimSpace(tool)
			if tool == "" {
				continue
			}
			toolNames = append(toolNames, tool)
		}
		normalized[serverIDText] = toolNames
	}

	body, err := json.Marshal(normalized)
	if err != nil {
		return fmt.Errorf("MCP服务工具配置编码失败: %w", err)
	}
	req.AllowMcpServerIdsTools = normalized
	req.AllowMcpServerIdsToolsRaw = string(body)
	return nil
}
