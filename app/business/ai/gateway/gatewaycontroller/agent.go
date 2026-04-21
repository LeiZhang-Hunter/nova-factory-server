package gatewaycontroller

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"strings"

	"nova-factory-server/app/business/admin/system/systemdao"
	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Agent 智能体控制器。
type Agent struct {
	service              gatewayservice.IAIAgentService
	modelProviderService aidatasetservice.IAiModelProviderService
	agentDao             gatewaydao.IAIAgentDao
	mcpServerDao         gatewaydao.IMCPServerDao
	llmDao               aidatasetdao.IAiLLMDao
	dictDataDao          systemdao.IDictDataDao
}

// NewAgent 创建智能体控制器。
func NewAgent(service gatewayservice.IAIAgentService, modelProviderService aidatasetservice.IAiModelProviderService, agentDao gatewaydao.IAIAgentDao, mcpServerDao gatewaydao.IMCPServerDao, llmDao aidatasetdao.IAiLLMDao, dictDataDao systemdao.IDictDataDao) *Agent {
	return &Agent{
		service:              service,
		modelProviderService: modelProviderService,
		agentDao:             agentDao,
		mcpServerDao:         mcpServerDao,
		llmDao:               llmDao,
		dictDataDao:          dictDataDao,
	}
}

// PrivateRoutes 注册智能体配置路由。
func (agent *Agent) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/agent/config")
	group.GET("/list", middlewares.HasPermission("ai:agent:config:list"), agent.List)
	group.GET("/query/:id", middlewares.HasPermission("ai:agent:config:query"), agent.GetByID)
	group.GET("/query/type/:type", middlewares.HasPermission("ai:agent:config:type"), agent.GetEnabledByType)
	group.POST("/set", middlewares.HasPermission("ai:agent:config:set"), agent.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("ai:agent:config:remove"), agent.Delete)
}

// List 获取智能体列表
// @Summary 获取智能体列表
// @Description 获取智能体列表
// @Tags 工业智能体/智能体管理
// @Param object query gatewaymodels.AIAgentQuery true "智能体查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/agent/config/list [get]
func (agent *Agent) List(c *gin.Context) {
	req := new(gatewaymodels.AIAgentQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := agent.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetByID 获取智能体详情
// @Summary 获取智能体详情
// @Description 根据ID获取智能体详情
// @Tags 工业智能体/智能体管理
// @Param id path int true "智能体ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/agent/config/query/{id} [get]
func (agent *Agent) GetByID(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := agent.service.GetByID(c, id)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetEnabledByType 获取指定类型下已启用的智能体详情
// @Summary 获取指定类型下已启用的智能体详情
// @Description 通过type查询enable=true的智能体配置
// @Tags 工业智能体/智能体管理
// @Param type path string true "智能体类型"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/agent/config/query/type/{type} [get]
func (agent *Agent) GetEnabledByType(c *gin.Context) {
	agentType := strings.TrimSpace(c.Param("type"))
	if agentType == "" {
		baizeContext.ParameterError(c)
		return
	}
	data, err := agent.service.GetEnabledByType(c, agentType)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set 保存智能体
// @Summary 保存智能体
// @Description 保存智能体，id为空时新增，不为空时修改
// @Tags 工业智能体/智能体管理
// @Param object body gatewaymodels.AIAgentUpsert true "智能体保存参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /ai/agent/config/set [post]
func (agent *Agent) Set(c *gin.Context) {
	req := new(gatewaymodels.AIAgentUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if err := agent.validateType(c, req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	if err := agent.validateSingleEnabledByType(c, req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	if err := agent.validateDefaultLLM(c, req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	if err := agent.prepareAllowMcpServerIdsTools(c, req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	var (
		data *gatewaymodels.AIAgent
		err  error
	)
	if req.ID != 0 {
		data, err = agent.service.Update(c, req)
	} else {
		data, err = agent.service.Create(c, req)
	}
	if err != nil {
		zap.L().Error("set agent error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

func (agent *Agent) prepareAllowMcpServerIdsTools(c *gin.Context, req *gatewaymodels.AIAgentUpsert) error {
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
		server, err := agent.mcpServerDao.GetByID(c, serverID)
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
		return fmt.Errorf("允许使用的MCP工具配置编码失败: %w", err)
	}
	req.AllowMcpServerIdsTools = normalized
	req.AllowMcpServerIdsToolsRaw = string(body)
	return nil
}

// Delete 删除智能体
// @Summary 删除智能体
// @Description 删除智能体
// @Tags 工业智能体/智能体管理
// @Param ids path string true "智能体ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /ai/agent/config/remove/{ids} [delete]
func (agent *Agent) Delete(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := agent.service.DeleteByIDs(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

func (agent *Agent) validateType(c *gin.Context, req *gatewaymodels.AIAgentUpsert) error {
	if req == nil {
		return errors.New("参数不能为空")
	}
	req.Type = strings.TrimSpace(req.Type)
	if req.Type == "" {
		return errors.New("type不能为空")
	}
	typeKey := strings.ToLower(req.Type)
	rows := agent.dictDataDao.SelectDictDataByType(c, "ai_agent_type")
	for _, row := range rows {
		if row == nil {
			continue
		}
		if strings.ToLower(strings.TrimSpace(row.DictValue)) == typeKey || strings.ToLower(strings.TrimSpace(row.DictLabel)) == typeKey {
			return nil
		}
	}
	return errors.New("type 不存在于 ai_agent_type 字典")
}

func (agent *Agent) validateSingleEnabledByType(c *gin.Context, req *gatewaymodels.AIAgentUpsert) error {
	if req == nil {
		return errors.New("参数不能为空")
	}
	if req.Enable == nil || !*req.Enable {
		return nil
	}
	typeKey := strings.ToLower(strings.TrimSpace(req.Type))
	if typeKey == "" {
		return nil
	}
	current, err := agent.agentDao.GetEnabledByType(c, req.Type)
	if err != nil {
		return err
	}
	if current == nil {
		return nil
	}
	if req.ID != 0 && current.ID == req.ID {
		return nil
	}
	if strings.ToLower(strings.TrimSpace(current.Type)) == typeKey {
		return errors.New("同一 type 只能启用一个智能体")
	}
	return nil
}

func (agent *Agent) validateDefaultLLM(c *gin.Context, req *gatewaymodels.AIAgentUpsert) error {
	if req == nil {
		return errors.New("参数不能为空")
	}
	req.DefaultLLMProviderID = strings.TrimSpace(req.DefaultLLMProviderID)
	req.DefaultLLMModelID = strings.TrimSpace(req.DefaultLLMModelID)
	if req.DefaultLLMProviderID == "" && req.DefaultLLMModelID == "" {
		return nil
	}
	if req.DefaultLLMProviderID == "" || req.DefaultLLMModelID == "" {
		return errors.New("defaultLlmProviderId 和 defaultLlmModelId 不能为空")
	}
	providers, err := agent.modelProviderService.ListWithLLM(c, &aidatasetmodels.SysAiModelProviderListReq{
		Name: req.DefaultLLMProviderID,
	})
	if err != nil {
		return err
	}
	providerKey := strings.ToLower(req.DefaultLLMProviderID)
	var matchedProvider *aidatasetmodels.SysAiModelProvider
	for _, provider := range providers.Rows {
		if provider == nil {
			continue
		}
		if strings.ToLower(provider.Name) == strings.ToLower(providerKey) || strconv.FormatInt(provider.ID, 10) == req.DefaultLLMProviderID {
			matchedProvider = provider
			break
		}
	}
	if matchedProvider == nil {
		return errors.New("defaultLlmProviderId 不存在")
	}
	llms, err := agent.llmDao.ListByFactory(c, matchedProvider.Name)
	if err != nil {
		return err
	}
	for _, llm := range llms {
		if llm == nil {
			continue
		}
		if llm.LlmName == req.DefaultLLMModelID || llm.ModelType == req.DefaultLLMModelID {
			return nil
		}
	}
	return errors.New("defaultLlmModelId 不存在")
}
