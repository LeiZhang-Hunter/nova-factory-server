package aiDataSetController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/business/ai/aiDataSetService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Prediction struct {
	service aiDataSetService.IAiPredictionService
}

func NewPrediction(service aiDataSetService.IAiPredictionService) *Prediction {
	return &Prediction{
		service: service,
	}
}

func (p *Prediction) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/prediction")
	group.POST("/list", middlewares.HasPermission("ai:prediction:list"), p.List)
	group.POST("/set", middlewares.HasPermission("ai:prediction:set"), p.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("ai:prediction:remove"), p.Remove)
	return
}

// List 智能预警列表
// @Summary 智能预警列表
// @Description 智能预警列表
// @Tags 工业智能体/智能预警
// @Param  object query aiDataSetModels.SysAiPredictionListReq true "智能预警列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/prediction/list [get]
func (p *Prediction) List(c *gin.Context) {
	req := new(aiDataSetModels.SysAiPredictionListReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := p.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// Set 设置智能预警
// @Summary 设置智能预警
// @Description 设置智能预警
// @Tags 工业智能体/智能预警
// @Param  object body aiDataSetModels.SetSysAiPrediction true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/prediction/set [post]
func (p *Prediction) Set(c *gin.Context) {
	req := new(aiDataSetModels.SetSysAiPrediction)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := p.service.Set(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// Remove 删除智能预警
// @Summary 删除智能预警
// @Description 删除智能预警
// @Tags 工业智能体/智能预警
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /ai/prediction/remove/{ids}  [delete]
func (p *Prediction) Remove(c *gin.Context) {
	ids := baizeContext.ParamStringArray(c, "ids")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		return
	}
	err := p.service.Remove(c, ids)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}

	baizeContext.Success(c)
}
