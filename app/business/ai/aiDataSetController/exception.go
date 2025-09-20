package aiDataSetController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/business/ai/aiDataSetService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Exception struct {
	service aiDataSetService.IAiPredictionExceptionService
}

func NewException(service aiDataSetService.IAiPredictionExceptionService) *Exception {
	return &Exception{
		service: service,
	}
}

func (e *Exception) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/exception")
	group.GET("/list", middlewares.HasPermission("ai:exception:list"), e.List)
	group.POST("/set", middlewares.HasPermission("ai:exception:set"), e.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("ai:exception:remove"), e.Remove)
	return
}

// List 异常预警列表
// @Summary 异常预警列表
// @Description 异常预警列表
// @Tags 工业智能体/异常预警
// @Param  object query aiDataSetModels.SysAiPredictionListReq true "异常预警参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/exception/list [get]
func (e *Exception) List(c *gin.Context) {
	req := new(aiDataSetModels.SysAiPredictionExceptionListReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		zap.L().Warn("参数错误", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	list, err := e.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// Set 设置异常预警
// @Summary 设置异常预警
// @Description 设置异常预警
// @Tags 工业智能体/异常预警
// @Param  object query aiDataSetModels.SysAiPredictionListReq true "设置异常预警参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/exception/set [post]
func (e *Exception) Set(c *gin.Context) {
	req := new(aiDataSetModels.SetSysAiPredictionException)
	err := c.ShouldBindJSON(req)
	if err != nil {
		zap.L().Warn("参数错误", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	list, err := e.service.Set(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// Remove 删除异常预警
// @Summary 删除异常预警
// @Description 删除异常预警
// @Tags 工业智能体/异常预警
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /ai/exception/remove/{ids}  [delete]
func (e *Exception) Remove(c *gin.Context) {
	ids := baizeContext.ParamStringArray(c, "ids")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		return
	}
	err := e.service.Remove(c, ids)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}

	baizeContext.Success(c)
}
