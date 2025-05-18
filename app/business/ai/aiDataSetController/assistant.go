package aiDataSetController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/utils/baizeContext"
)

// CreateAssistant 创建助理
// @Summary 创建助理
// @Description 创建助理
// @Tags 工业智能体/助理管理
// @Param  object body aiDataSetModels.CreateAssistantRequest true "创建助理表请求参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/assistant/create [post]
func (d *Dataset) CreateAssistant(c *gin.Context) {
	req := new(aiDataSetModels.CreateAssistantRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		zap.L().Error("添加助理失败", zap.Error(err))
		baizeContext.Waring(c, "添加助理失败")
		return
	}
	assistant, err := d.assistantService.AddAssistant(c, req)
	if err != nil {
		zap.L().Error("添加助理失败", zap.Error(err))
		baizeContext.SuccessData(c, "添加助理失败")
		return
	}
	baizeContext.SuccessData(c, assistant.Data)
	return
}

// UpdateAssistant 更新助理
// @Summary 更新助理
// @Description 更新助理
// @Tags 工业智能体/助理管理
// @Param  object body aiDataSetModels.UpdateAssistantRequest true "更新助理表请求参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/assistant/update [put]
func (d *Dataset) UpdateAssistant(c *gin.Context) {
	req := new(aiDataSetModels.UpdateAssistantRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		zap.L().Error("更新助理失败", zap.Error(err))
		baizeContext.Waring(c, "更新助理失败")
		return
	}
	if req.ChatId == "" {
		baizeContext.Waring(c, "请输入聊天id")
		return
	}
	assistant, err := d.assistantService.UpdateAssistant(c, req)
	if err != nil {
		zap.L().Error("更新助理失败", zap.Error(err))
		baizeContext.SuccessData(c, "更新助理失败")
		return
	}
	baizeContext.SuccessData(c, assistant)
	return
}

// RemoveAssistant 删除助理
// @Summary 删除助理
// @Description 删除助理
// @Tags 工业智能体/助理管理
// @Param  assistantIds path []string true "assistantIds"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/assistant/remove/{assistantIds} [delete]
func (d *Dataset) RemoveAssistant(c *gin.Context) {
	assistantIds := baizeContext.ParamStringArray(c, "assistantIds")
	err := d.assistantService.DeleteAssistant(c, assistantIds)
	if err != nil {
		baizeContext.Waring(c, "助理删除失败")
		return
	}
	baizeContext.Success(c)
}

// ListAssistant 读取助理列表
// @Summary 读取助理列表
// @Description 读取助理列表
// @Tags 工业智能体/助理管理
// @Param  object query aiDataSetModels.GetAssistantRequest true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/assistant/list [get]
func (d *Dataset) ListAssistant(c *gin.Context) {
	req := new(aiDataSetModels.GetAssistantRequest)
	err := c.ShouldBind(req)
	if err != nil {
		zap.L().Error("更新助理失败", zap.Error(err))
		baizeContext.Waring(c, "更新助理失败")
		return
	}
	assistants, err := d.assistantService.ListAssistant(c, req)
	if err != nil {
		baizeContext.Waring(c, "读取助理列表失败")
		return
	}
	baizeContext.SuccessData(c, assistants.Data)
}
