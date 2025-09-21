package aiDataSetController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/business/ai/aiDataSetService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Control struct {
	service aiDataSetService.IAiPredictionControlService
}

func NewControl(service aiDataSetService.IAiPredictionControlService) *Control {
	return &Control{
		service: service,
	}
}

func (i *Control) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/control")
	group.GET("/list", middlewares.HasPermission("ai:control:list"), i.List)
	group.POST("/set", middlewares.HasPermission("ai:control:set"), i.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("ai:control:remove"), i.Remove)
	return
}

// List 预测控制列表
// @Summary 预测控制列表
// @Description 预测控制列表
// @Tags 工业智能体/预测控制
// @Param  object query aiDataSetModels.SysAiPredictionControlListReq true "智能控制列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/control/list [get]
func (i *Control) List(c *gin.Context) {
	req := new(aiDataSetModels.SysAiPredictionControlListReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		zap.L().Warn("参数错误", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	list, err := i.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// Set 设置智能控制
// @Summary 设置智能控制
// @Description 设置智能控制
// @Tags 工业智能体/预测控制
// @Param  object body aiDataSetModels.SetSysAiPredictionControl true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/prediction/set [post]
func (i *Control) Set(c *gin.Context) {
	req := new(aiDataSetModels.SetSysAiPredictionControl)
	err := c.ShouldBindJSON(req)
	if err != nil {
		zap.L().Warn("参数错误", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	list, err := i.service.Set(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// Remove 删除预测控制
// @Summary 删除预测控制
// @Description 删除预测控制
// @Tags 工业智能体/预测控制
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /ai/prediction/remove/{ids}  [delete]
func (i *Control) Remove(c *gin.Context) {
	ids := baizeContext.ParamStringArray(c, "ids")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		return
	}
	err := i.service.Remove(c, ids)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}

	baizeContext.Success(c)
}
