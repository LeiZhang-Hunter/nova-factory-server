package deviceController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type DeviceCheckPlan struct {
	service deviceService.IDeviceCheckPlanService
}

func NewDeviceCheckPlan(service deviceService.IDeviceCheckPlanService) *DeviceCheckPlan {
	return &DeviceCheckPlan{
		service: service,
	}
}

func (dc *DeviceCheckPlan) PrivateRoutes(router *gin.RouterGroup) {
	subject := router.Group("/asset/device/checkPlan")
	subject.GET("/list", middlewares.HasPermission("asset:device:checkPlan:list"), dc.List)               // 点检保养项目列表
	subject.POST("/set", middlewares.HasPermission("asset:device:checkPlan:set"), dc.Set)                 // 设置点检保养项目
	subject.DELETE("/remove/:ids", middlewares.HasPermission("asset:device:checkPlan:remove"), dc.Remove) //删除点检保养项目
}

// List 点检保养计划列表
// @Summary 点检保养计划列表
// @Description 点检保养计划列表
// @Tags 设备管理/点检保养计划
// @Param  object query deviceModels.SysDeviceCheckPlanReq true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/device/checkPlan/list [get]
func (dc *DeviceCheckPlan) List(c *gin.Context) {
	req := new(deviceModels.SysDeviceCheckPlanReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	list, err := dc.service.List(c, req)
	if err != nil {
		zap.L().Error("读取点检保养计划失败", zap.Error(err))
		baizeContext.Waring(c, "读取点检保养计划失败")
		return
	}
	baizeContext.SuccessData(c, list)
}

// Set 设置点检保养计划
// @Summary 设置点检保养计划
// @Description 设置点检保养计划
// @Tags 设备管理/点检保养计划
// @Param  object body deviceModels.SysDeviceCheckPlanVO true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/device/checkPlan/set [post]
func (dc *DeviceCheckPlan) Set(c *gin.Context) {
	vo := new(deviceModels.SysDeviceCheckPlanVO)
	err := c.ShouldBindJSON(vo)
	if err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	set, err := dc.service.Set(c, vo)
	if err != nil {
		zap.L().Error("设置点检保养计划失败", zap.Error(err))
		baizeContext.Waring(c, "设置点检保养计划失败")
		return
	}
	baizeContext.SuccessData(c, set)
}

// Remove 删除点检保养计划
// @Summary 删除点检保养计划
// @Description 删除点检保养计划
// @Tags 设备管理/点检保养计划
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /asset/device/checkPlan/remove/{ids}  [delete]
func (dc *DeviceCheckPlan) Remove(c *gin.Context) {
	ids := baizeContext.ParamStringArray(c, "ids")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		baizeContext.ParameterError(c)
		return
	}
	err := dc.service.Remove(c, ids)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}

	baizeContext.Success(c)
}
