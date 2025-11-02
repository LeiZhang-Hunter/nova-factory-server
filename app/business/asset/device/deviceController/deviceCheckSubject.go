package deviceController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type DeviceCheckSubject struct {
	service deviceService.IDeviceCheckSubjectService
}

func NewDeviceCheckSubject(service deviceService.IDeviceCheckSubjectService) *DeviceCheckSubject {
	return &DeviceCheckSubject{
		service: service,
	}
}

func (d *DeviceCheckSubject) PrivateRoutes(router *gin.RouterGroup) {
	subject := router.Group("/asset/device/check/subject")
	subject.GET("/list", middlewares.HasPermission("asset:device:check:subject:list"), d.List)               // 点检保养项目列表
	subject.POST("/set", middlewares.HasPermission("asset:device:check:subject:set"), d.Set)                 // 设置点检保养项目
	subject.DELETE("/remove/:ids", middlewares.HasPermission("asset:device:check:subject:remove"), d.Remove) //删除点检保养项目
}

// Set 设置点检项目
// @Summary 设置点检项目
// @Description 设置点检项目
// @Tags 设备管理/点检保养计划/点检项目
// @Param  object body deviceModels.SysDeviceCheckSubjectVO true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/device/check/subject/set [post]
func (d *DeviceCheckSubject) Set(c *gin.Context) {
	vo := new(deviceModels.SysDeviceCheckSubjectVO)
	err := c.ShouldBindJSON(vo)
	if err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	set, err := d.service.Set(c, vo)
	if err != nil {
		zap.L().Error("设置点检保养计划失败", zap.Error(err))
		baizeContext.Waring(c, "设置点检保养计划失败")
		return
	}
	baizeContext.SuccessData(c, set)
}

// Remove 删除点检项目
// @Summary 删除点检项目
// @Description 删除点检项目
// @Tags 设备管理/点检保养计划/点检项目
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /asset/device/check/subject/remove/{ids}  [delete]
func (d *DeviceCheckSubject) Remove(c *gin.Context) {
	ids := baizeContext.ParamStringArray(c, "ids")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		baizeContext.ParameterError(c)
		return
	}
	err := d.service.Remove(c, ids)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}

	baizeContext.Success(c)
}

// List 点检项目列表
// @Summary 点检项目列表
// @Description 点检项目列表
// @Tags 设备管理/点检保养计划/点检项目
// @Param  object query deviceModels.SysDeviceCheckSubjectReq true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/device/check/subject/list [get]
func (d *DeviceCheckSubject) List(c *gin.Context) {
	req := new(deviceModels.SysDeviceCheckSubjectReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	list, err := d.service.List(c, req)
	if err != nil {
		zap.L().Error("读取点检保养计划失败", zap.Error(err))
		baizeContext.Waring(c, "读取点检保养计划失败")
		return
	}
	baizeContext.SuccessData(c, list)
}
