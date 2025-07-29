package alertController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/business/alert/alertService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type AlertTemplate struct {
	service alertService.AlertTemplateService
}

func NewAlertTemplate(service alertService.AlertTemplateService) *AlertTemplate {
	return &AlertTemplate{
		service: service,
	}
}

func (ac *AlertTemplate) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/alert/template")
	group.GET("/list", middlewares.HasPermission("alert:template:list"), ac.List)               // 物料列表
	group.POST("/set", middlewares.HasPermission("alert:template:set"), ac.Set)                 // 设置物料信息
	group.DELETE("/remove/:ids", middlewares.HasPermission("alert:template:remove"), ac.Remove) //删除物料
}

// List 告警模板列表
// @Summary 告警模板列表
// @Description 告警模板列表
// @Tags 告警管理/告警模板管理
// @Param  object query alertModels.SysAlertSinkTemplateReq true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /alert/template/list [get]
func (ac *AlertTemplate) List(c *gin.Context) {
	req := new(alertModels.SysAlertSinkTemplateReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := ac.service.List(c, req)
	if err != nil {
		zap.L().Error("get template list error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
	return
}

// Set 设置告警发送配置
// @Summary 设置告警发送配置
// @Description 设置告警发送配置
// @Tags 告警管理/告警模板管理
// @Param  object body alertModels.SetSysAlertSinkTemplate true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /alert/template/set [post]
func (ac *AlertTemplate) Set(c *gin.Context) {
	info := new(alertModels.SetSysAlertSinkTemplate)
	err := c.ShouldBindJSON(info)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if info.ID == 0 {
		value, err := ac.service.Create(c, info)
		if err != nil {
			zap.L().Error("create template error", zap.Error(err))
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, alertModels.FromSysAlertSinkTemplateToSet(value))
	} else {
		value, err := ac.service.Update(c, info)
		if err != nil {
			zap.L().Error("create template error", zap.Error(err))
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, alertModels.FromSysAlertSinkTemplateToSet(value))
	}
}

// Remove 删除告警发送配置
// @Summary 删除告警发送配置
// @Description 删除告警发送配置
// @Tags 告警管理/告警模板管理
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /alert/template/remove/{ids}  [delete]
func (ac *AlertTemplate) Remove(c *gin.Context) {
	ids := baizeContext.ParamStringArray(c, "ids")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		return
	}
	err := ac.service.Remove(c, ids)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}

	baizeContext.Success(c)
}
