package alertController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/business/alert/alertService"
	"nova-factory-server/app/business/system/systemService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/time"
)

type AlertAction struct {
	service alertService.AlertActionService
}

func NewAlertAction(service alertService.AlertActionService, dds systemService.IDictDataService) *AlertAction {
	return &AlertAction{
		service: service,
	}
}

func (a *AlertAction) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/alert/action")
	group.GET("/list", middlewares.HasPermission("alert:action:list"), a.List)               // 告警规则列表
	group.POST("/set", middlewares.HasPermission("alert:action:set"), a.Set)                 // 设置物料信息
	group.DELETE("/remove/:ids", middlewares.HasPermission("alert:action:remove"), a.Remove) //删除物料
	return
}

// List 告警动作列表
// @Summary 告警动作列表
// @Description 告警动作列表
// @Tags 告警管理/告警处理管理
// @Param  object query alertModels.SysAlertActionListReq true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /alert/action/list [get]
func (a *AlertAction) List(c *gin.Context) {
	req := new(alertModels.SysAlertActionListReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := a.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// Set 设置告警发送配置
// @Summary 设置告警发送配置
// @Description 设置告警发送配置
// @Tags 告警管理/告警处理管理
// @Param  object body alertModels.SetAlertAction true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /alert/action/set [post]
func (a *AlertAction) Set(c *gin.Context) {
	info := new(alertModels.SetAlertAction)
	err := c.ShouldBindJSON(info)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	//校验用户通知
	if len(info.ApiNotifyList) == 0 && len(info.UserNotifyList) == 0 {
		baizeContext.Waring(c, "通知插件至少设置一项")
		return
	}

	// 校验用户通知
	if len(info.UserNotifyList) > 0 {
		for k, userNotify := range info.UserNotifyList {
			// 校验周期
			var checkWeeks map[string]int = map[string]int{
				"1": 0,
				"2": 0,
				"3": 0,
				"4": 0,
				"5": 0,
				"6": 0,
				"7": 0,
			}
			if len(userNotify.Period) == 0 {
				baizeContext.Waring(c, "用户通知通知周期不能为空")
				return
			}
			if len(userNotify.Receiver) == 0 {
				baizeContext.Waring(c, "用户通知通知人不能为空")
				return
			}

			if len(userNotify.TimeRange) == 0 {
				baizeContext.Waring(c, "用户通知通知时段不能为空")
				return
			}

			if len(userNotify.TimeRange) != 2 {
				baizeContext.Waring(c, "用户通知通知时段格式错误")
				return
			}

			value, err := time.TimeToString(userNotify.TimeRange[0])
			if err != nil {
				zap.L().Error("TimeToString error", zap.Error(err))
				baizeContext.Waring(c, "用户通知通知时段格式错误")
				return
			}
			info.UserNotifyList[k].TimeStart = uint64(value)

			valueEnd, err := time.TimeToString(userNotify.TimeRange[1])
			if err != nil {
				zap.L().Error("TimeToString error", zap.Error(err))
				baizeContext.Waring(c, "用户通知通知时段格式错误")
				return
			}
			info.UserNotifyList[k].TimeEnd = uint64(valueEnd)
			for k, week := range userNotify.Period {
				value, ok := checkWeeks[week]
				if !ok {
					baizeContext.Waring(c, "用户通知通知周期错误")
					return
				}
				if value != 0 {
					baizeContext.Waring(c, fmt.Sprintf("执行日期格式错误，存在重复日期:%d", k))
				}
				checkWeeks[week]++
			}
		}
	}

	// 校验api通知
	if len(info.ApiNotifyList) > 0 {
		for k, apiNotify := range info.ApiNotifyList {
			// 校验周期
			var checkWeeks map[string]int = map[string]int{
				"1": 0,
				"2": 0,
				"3": 0,
				"4": 0,
				"5": 0,
				"6": 0,
				"7": 0,
			}
			if len(apiNotify.Period) == 0 {
				baizeContext.Waring(c, "接口回调通知周期不能为空")
				return
			}
			if len(apiNotify.TimeRange) == 0 {
				baizeContext.Waring(c, "用户通知通知时段不能为空")
				return
			}
			if len(apiNotify.TimeRange) != 2 {
				baizeContext.Waring(c, "用户通知通知时段格式错误")
				return
			}

			value, err := time.TimeToString(apiNotify.TimeRange[0])
			if err != nil {
				zap.L().Error("TimeToString error", zap.Error(err))
				baizeContext.Waring(c, "用户通知通知时段格式错误")
				return
			}
			info.ApiNotifyList[k].TimeStart = uint64(value)

			valueEnd, err := time.TimeToString(apiNotify.TimeRange[1])
			if err != nil {
				zap.L().Error("TimeToString error", zap.Error(err))
				baizeContext.Waring(c, "用户通知通知时段格式错误")
				return
			}
			info.ApiNotifyList[k].TimeStart = uint64(valueEnd)

			for k, week := range apiNotify.Period {
				value, ok := checkWeeks[week]
				if !ok {
					baizeContext.Waring(c, "用户通知通知周期错误")
					return
				}
				if value != 0 {
					baizeContext.Waring(c, fmt.Sprintf("执行日期格式错误，存在重复日期:%d", k))
				}
				checkWeeks[week]++
			}
		}
	}

	if info.ID == 0 {
		value, err := a.service.Set(c, info)
		if err != nil {
			zap.L().Error("create action error", zap.Error(err))
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, value)
	} else {
		value, err := a.service.Set(c, info)
		if err != nil {
			zap.L().Error("create action error", zap.Error(err))
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, value)
	}
}

// Remove 删除告警发送配置
// @Summary 删除告警发送配置
// @Description 删除告警发送配置
// @Tags 告警管理/告警处理管理
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /alert/action/remove/{ids}  [delete]
func (a *AlertAction) Remove(c *gin.Context) {
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
