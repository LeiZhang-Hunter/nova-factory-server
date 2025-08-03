package alertController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/business/alert/alertService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type AlertAiReason struct {
	service alertService.AlertAiReasonService
}

func NewAlertAiReason(service alertService.AlertAiReasonService) *AlertAiReason {
	return &AlertAiReason{
		service: service,
	}
}

func (a *AlertAiReason) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/alert/reason")
	group.GET("/list", middlewares.HasPermission("alert:reason:list"), a.List)               // 告警规则列表
	group.POST("/set", middlewares.HasPermission("alert:reason:set"), a.Set)                 // 设置物料信息
	group.DELETE("/remove/:ids", middlewares.HasPermission("alert:reason:remove"), a.Remove) //删除物料
	return
}

// List 告警AI推理列表
// @Summary 告警AI推理列表
// @Description 告警AI推理列表
// @Tags 告警管理/告警AI推理管理
// @Param  object query alertModels.SysAlertAiReasonReq true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /alert/reason/list [get]
func (a *AlertAiReason) List(c *gin.Context) {
	req := new(alertModels.SysAlertAiReasonReq)
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

// Set 设置告警AI推理发送配置
// @Summary 设置告警AI推理发送配置
// @Description 设置告警AI推理发送配置
// @Tags 告警AI推理管理/告警AI推理管理
// @Param  object body alertModels.SetAlertAiReason true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /alert/reason/set [post]
func (a *AlertAiReason) Set(c *gin.Context) {
	info := new(alertModels.SetAlertAiReason)
	err := c.ShouldBindJSON(info)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	value, err := a.service.Set(c, info)
	baizeContext.SuccessData(c, value)

}

// Remove 删除告警AI推理发送配置
// @Summary 删除告警AI推理发送配置
// @Description 删除告警AI推理发送配置
// @Tags 告警AI推理管理/告警AI推理管理
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /alert/reason/remove/{ids}  [delete]
func (a *AlertAiReason) Remove(c *gin.Context) {
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
