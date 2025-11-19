package alertController

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/business/alert/alertService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Alert struct {
	service alertService.AlertRuleService
	action  *AlertAction
}

func NewAlert(service alertService.AlertRuleService) *Alert {
	return &Alert{
		service: service,
	}
}

func (a *Alert) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/alert/rule")
	group.GET("/list", middlewares.HasPermission("alert:rule:list"), a.List)               // 告警规则列表
	group.POST("/set", middlewares.HasPermission("alert:rule:set"), a.Set)                 // 设置物料信息
	group.DELETE("/remove/:ids", middlewares.HasPermission("alert:rule:remove"), a.Remove) //删除物料
}

// List 告警规则列表
// @Summary 告警规则列表
// @Description 告警规则列表
// @Tags 告警管理/告警规则管理
// @Param  object query alertModels.SysAlertListReq true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /alert/rule/list [get]
func (a *Alert) List(c *gin.Context) {
	req := new(alertModels.SysAlertListReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := a.service.List(c, req)
	if err != nil {
		zap.L().Error("get rule list error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
	return
}

// Set 设置告警规则
// @Summary 设置告警规则
// @Description 设置告警规则
// @Tags 告警管理/告警规则管理
// @Param object body alertModels.SetSysAlert true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /alert/rule/set [post]
func (a *Alert) Set(c *gin.Context) {
	info := new(alertModels.SetSysAlert)
	err := c.ShouldBindJSON(info)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	// 校验rule是否合法
	if len(info.Advanced.Rules) != 0 {
		for _, rule := range info.Advanced.Rules {
			for _, group := range rule.Groups {
				if group.Key == "" {
					baizeContext.Waring(c, "数据id不能为空")
					return
				}
				if group.Name == "" {
					baizeContext.Waring(c, "数据id名称不能为空")
					return
				}
				if group.Operator == "" {
					baizeContext.Waring(c, "条件操作符号不能为空")
					return
				}
				if group.OperatorName == "" {
					baizeContext.Waring(c, "条件操作符号名称不能为空")
					return
				}
				if group.Value == "" {
					baizeContext.Waring(c, "规则数值不能是空")
					return
				}
			}
		}

	}

	if info.ID == 0 {
		value, err := a.service.Create(c, info)
		if err != nil {
			zap.L().Error("create rule error", zap.Error(err))
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, alertModels.FromSysAlertToSetData(value))
	} else {
		// 查询是否有开启策略，只能开启一个
		open, err := a.service.FindOpen(c, info.GatewayID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				zap.L().Error("find rule error", zap.Error(err))
				baizeContext.Waring(c, "读取开启数据失败")
				return
			}

		}

		if open != nil {
			if open.ID != info.ID {
				baizeContext.Waring(c, "告警策略只能开启一个")
				return
			}
		}

		value, err := a.service.Update(c, info)
		if err != nil {
			zap.L().Error("update rule error", zap.Error(err))
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, alertModels.FromSysAlertToSetData(value))
	}
}

// Remove 删除告警规则
// @Summary 删除告警规则
// @Description 删除告警规则
// @Tags 告警管理/告警规则管理
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /alert/rule/remove/{ids}  [delete]
func (a *Alert) Remove(c *gin.Context) {
	ids := baizeContext.ParamStringArray(c, "ids")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		return
	}
	err := a.service.Remove(c, ids)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}

	baizeContext.Success(c)
}

// Change 改变告警配置
// @Summary 改变告警配置
// @Description 改变告警配置
// @Tags 告警管理/告警规则管理
// @Param  object body alertModels.ChangeSysAlert true "助理列表参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /alert/rule/change  [post]
func (a *Alert) Change(c *gin.Context) {
	info := new(alertModels.ChangeSysAlert)
	err := c.ShouldBindJSON(info)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	err = a.service.Change(c, info)
	if err != nil {
		zap.L().Error("change rule error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
