package aiDataSetController

import (
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"
	"nova-factory-server/app/business/ai/agent/aiDataSetService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Model struct {
	service        aiDataSetService.IAiModelProviderService
	settingService aiDataSetService.IAiLLMSettingService
	userLLMService aiDataSetService.IAiUserLLMService
}

func NewModel(service aiDataSetService.IAiModelProviderService, settingService aiDataSetService.IAiLLMSettingService, userLLMService aiDataSetService.IAiUserLLMService) *Model {
	return &Model{
		service:        service,
		settingService: settingService,
		userLLMService: userLLMService,
	}
}

func (m *Model) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/model")
	group.GET("/provider/list", middlewares.HasPermission("ai:model:provider:list"), m.ProviderList)
	group.GET("/provider/setting/get", middlewares.HasPermission("ai:model:provider:setting:get"), m.GetSetting)
	group.POST("/provider/setting/set", middlewares.HasPermission("ai:model:provider:setting:set"), m.SetSetting)
	group.GET("/provider/global/get", middlewares.HasPermission("ai:model:provider:global:get"), m.GetGlobalModel)
	group.POST("/provider/global/set", middlewares.HasPermission("ai:model:provider:global:set"), m.SetGlobalModel)
}

// ProviderList 模型供应商列表
// @Summary 模型供应商列表
// @Description 读取模型供应商及其下级LLM列表
// @Tags 工业智能体/模型配置
// @Param  object query aiDataSetModels.SysAiModelProviderListReq true "模型供应商列表参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/model/provider/list [get]
func (m *Model) ProviderList(c *gin.Context) {
	req := new(aiDataSetModels.SysAiModelProviderListReq)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := m.service.ListWithLLM(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// SetSetting 设置模型配置
// @Summary 设置模型配置
// @Description 新增或修改模型配置
// @Tags 工业智能体/模型配置
// @Param  object body aiDataSetModels.SetSysAiLLMSetting true "模型配置参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /ai/model/provider/setting/set [post]
func (m *Model) SetSetting(c *gin.Context) {
	req := new(aiDataSetModels.SetSysAiLLMSetting)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := m.settingService.Set(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetSetting 读取模型配置
// @Summary 读取模型配置
// @Description 根据id读取模型配置，未传id时读取当前部门最近更新配置
// @Tags 工业智能体/模型配置
// @Param  object query aiDataSetModels.GetSysAiLLMSettingReq true "模型配置读取参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "读取成功"
// @Router /ai/model/provider/setting/get [get]
func (m *Model) GetSetting(c *gin.Context) {
	req := new(aiDataSetModels.GetSysAiLLMSettingReq)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := m.settingService.Get(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// SetGlobalModel 设置用户模型
// @Summary 设置用户模型
// @Description 新增或修改用户模型配置，user_id为0表示全局设置
// @Tags 工业智能体/模型配置
// @Param  object body aiDataSetModels.SetSysUserLLM true "全局模型配置参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /ai/model/provider/global/set [post]
func (m *Model) SetGlobalModel(c *gin.Context) {
	req := new(aiDataSetModels.SetSysUserLLM)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := m.userLLMService.Set(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetGlobalModel 读取用户模型
// @Summary 读取用户模型
// @Description 读取SetGlobalModel保存的用户模型配置
// @Tags 工业智能体/模型配置
// @Param  object query aiDataSetModels.GetSysUserLLMReq true "用户模型读取参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "读取成功"
// @Router /ai/model/provider/global/get [get]
func (m *Model) GetGlobalModel(c *gin.Context) {
	req := new(aiDataSetModels.GetSysUserLLMReq)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := m.userLLMService.Get(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}
