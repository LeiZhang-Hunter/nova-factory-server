package craftRouteController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type ProcessContext struct {
	processContextService craftRouteService.ICraftProcessContextService
}

func NewProcessContext(processContextService craftRouteService.ICraftProcessContextService) *ProcessContext {
	return &ProcessContext{
		processContextService: processContextService,
	}
}

func (p *ProcessContext) PrivateRoutes(router *gin.RouterGroup) {
	routers := router.Group("/craft/process/context")
	routers.GET("/list", middlewares.HasPermission("craft:process:context"), p.GetProcessContextList)                             // 工序内容列表
	routers.POST("/set", middlewares.HasPermission("craft:process:context:set"), p.SetProcessContextList)                         // 设置工序内容
	routers.DELETE("/remove/:context_ids", middlewares.HasPermission("craft:process:context:remove"), p.RemoveProcessContextList) //移除工序内容
}

// GetProcessContextList 工序内容列表
// @Summary 工序内容列表
// @Description 工序内容列表
// @Tags 工艺管理/工序内容管理
// @Param  object query craftRouteModels.SysProProcessContextListReq true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/process/context/list [get]
func (p *ProcessContext) GetProcessContextList(c *gin.Context) {
	req := new(craftRouteModels.SysProProcessContextListReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := p.processContextService.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// SetProcessContextList 设置工序内容
// @Summary 设置工序内容
// @Description 设置工序内容
// @Tags 工艺管理/工序内容管理
// @Param  object body craftRouteModels.SysProSetProcessContent true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/process/context/set [post]
func (p *ProcessContext) SetProcessContextList(c *gin.Context) {
	req := new(craftRouteModels.SysProSetProcessContent)
	err := c.ShouldBindJSON(req)
	if err != nil {
		zap.L().Error(err.Error())
		baizeContext.ParameterError(c)
		return
	}

	if req.ControlRules == nil {
		baizeContext.Waring(c, "控制算法不能为空")
		return
	}

	if req.ControlType != "pid" && req.ControlType != "threshold" && req.ControlType != "mpc" {
		baizeContext.Waring(c, "控制算法暂时不支持")
		return
	}

	if req.ControlType == "threshold" {
		if req.ControlRules.TriggerRules == nil {
			baizeContext.Waring(c, "阈值算法规则不能为空")
			return
		}

		if len(req.ControlRules.TriggerRules.Cases) == 0 {
			baizeContext.Waring(c, "阈值算法条件判断不能为空")
			return
		}

		for _, rule := range req.ControlRules.TriggerRules.Cases {
			if rule.Connector == "" {
				rule.Connector = "and"
			}
		}

		if len(req.ControlRules.TriggerRules.Actions) == 0 {
			baizeContext.Waring(c, "控制动作不能为空")
			return
		}

		for _, v := range req.ControlRules.TriggerRules.Actions {
			if v.DeviceId == "" {
				baizeContext.Waring(c, "调节设备不能为空")
				return
			}

			if v.DataId == "" {
				baizeContext.Waring(c, "调节数据不能为空")
				return
			}

			if v.ControlMode != "concurrent_control" && v.ControlMode != "delay_control" {
				baizeContext.Waring(c, "控制动作不能为空")
				return
			}

		}
	}

	if req.ControlType == "pid" {
		if req.ControlRules.PidRules == nil {
			baizeContext.Waring(c, "pid算法规则不能为空")
			return
		}
		if req.ControlRules.PidRules.DeviceId == "" {
			baizeContext.Waring(c, "请选择设备")
			return
		}
		if req.ControlRules.PidRules.DataId == "" {
			baizeContext.Waring(c, "请选择数据id")
			return
		}
		if req.ControlRules.PidRules.Proportional == 0 {
			baizeContext.Waring(c, "PID规则百分比增益不能为空")
			return
		}

		if len(req.ControlRules.PidRules.Actions) == 0 {
			baizeContext.Waring(c, "控制动作不能为空")
			return
		}

		for _, v := range req.ControlRules.PidRules.Actions {
			if v.DeviceId == "" {
				baizeContext.Waring(c, "调节设备不能为空")
				return
			}

			if v.DataId == "" {
				baizeContext.Waring(c, "调节数据不能为空")
				return
			}

			if v.ControlMode != "concurrent_control" && v.ControlMode != "delay_control" {
				baizeContext.Waring(c, "控制动作不能为空")
				return
			}

		}
	}

	if req.ControlType == "mpc" {
		if req.ControlRules.PredictRules == nil {
			baizeContext.Waring(c, "阈值算法规则不能为空")
			return
		}

		if len(req.ControlRules.PredictRules.Cases) == 0 {
			baizeContext.Waring(c, "阈值算法条件判断不能为空")
			return
		}

		for _, rule := range req.ControlRules.PredictRules.Cases {
			if rule.Connector == "" {
				rule.Connector = "and"
			}
		}

		if len(req.ControlRules.PredictRules.Actions) == 0 {
			baizeContext.Waring(c, "控制动作不能为空")
			return
		}

		for _, v := range req.ControlRules.PredictRules.Actions {
			if v.DeviceId == "" {
				baizeContext.Waring(c, "调节设备不能为空")
				return
			}

			if v.DataId == "" {
				baizeContext.Waring(c, "调节数据不能为空")
				return
			}

			if v.ControlMode != "concurrent_control" && v.ControlMode != "delay_control" {
				baizeContext.Waring(c, "控制动作不能为空")
				return
			}

		}

		if req.ControlRules.PredictRules.Interval <= 0 {
			baizeContext.Waring(c, "时间间隔格式错误")
			return
		}

		if req.ControlRules.PredictRules.PredictLength <= 0 {
			baizeContext.Waring(c, "预测长度错误")
			return
		}

		if req.ControlRules.PredictRules.Threshold <= 0 {
			baizeContext.Waring(c, "阀值错误")
			return
		}

		if req.ControlRules.PredictRules.Model == "" {
			baizeContext.Waring(c, "模型不能为空")
			return
		}

		if req.ControlRules.PredictRules.AggFunction == "" {
			baizeContext.Waring(c, "聚合函数不能为空")
			return
		}
	}

	if req.ContentID == 0 {
		ret, err := p.processContextService.Add(c, req)
		if err != nil {
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, ret)
	} else {
		ret, err := p.processContextService.Update(c, req)
		if err != nil {
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, ret)
	}
}

// RemoveProcessContextList 移除工序内容
// @Summary 移除工序内容
// @Description 移除工序内容
// @Tags 工艺管理/工序内容管理
// @Param  context_ids path string true "context_ids"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/process/context/remove/{context_ids} [delete]
func (p *ProcessContext) RemoveProcessContextList(c *gin.Context) {
	contextIds := baizeContext.ParamStringArray(c, "context_ids")
	if len(contextIds) == 0 {
		baizeContext.Waring(c, "请选择供需id")
		return
	}
	err := p.processContextService.Remove(c, contextIds)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}
	baizeContext.Success(c)
}
