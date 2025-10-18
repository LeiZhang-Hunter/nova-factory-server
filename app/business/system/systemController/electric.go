package systemController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/business/system/systemModels"
	"nova-factory-server/app/business/system/systemService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Electric struct {
	service       systemService.IDeviceElectricService
	deviceService deviceService.IDeviceService
}

// NewElectric  设备电流配置，通过电流判断设备到底是 运行负载中 还是停机 还是关闭
func NewElectric(service systemService.IDeviceElectricService, deviceService deviceService.IDeviceService) *Electric {
	return &Electric{
		service:       service,
		deviceService: deviceService,
	}
}

func (e *Electric) PrivateRoutes(router *gin.RouterGroup) {
	ele := router.Group("/system/electric")
	ele.GET("/list", middlewares.HasPermission("system:electric:list"), e.List)
	ele.POST("/set", middlewares.HasPermissions([]string{"system:electric:set"}), e.Set)
	ele.DELETE("/remove/:ids", middlewares.HasPermission("system:electric:remove"), e.Remove) //删除物料
}

// Set 设置设备电流配置
// @Summary 设置设备电流配置
// @Description 设置设备电流配置
// @Tags 系统管理/设备配置
// @Param  object body systemModels.SysDeviceElectricSettingVO true "设备电流配置参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /system/electric/set [post]
func (e *Electric) Set(c *gin.Context) {
	req := new(systemModels.SysDeviceElectricSettingVO)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.Expression == nil {
		baizeContext.Waring(c, "设备已经配置电流策略不能是空的")
		return
	}
	if len(req.Expression.Rules) == 0 {
		baizeContext.Waring(c, "设备已经配置电流策略不能是空的")
		return
	}
	info, err := e.deviceService.GetById(c, req.DeviceID)
	if err != nil {
		baizeContext.Waring(c, "读取设备数据错误")
		return
	}
	if info == nil {
		baizeContext.Waring(c, "设备不存在")
		return
	}
	if req.ID == 0 {
		electricInfo, err := e.service.GetByDeviceId(c, req.DeviceID)
		if err != nil {
			baizeContext.Waring(c, "设置错误")
			return
		}
		if electricInfo != nil {
			baizeContext.Waring(c, "设备已经配置电流策略不能重复设置")
			return
		}
	} else {
		electricInfo, err := e.service.GetByNoDeviceId(c, req.ID, req.DeviceID)
		if err != nil {
			baizeContext.Waring(c, "设置错误")
			return
		}
		if electricInfo != nil {
			baizeContext.Waring(c, "设备已经配置电流策略不能重复设置")
			return
		}
	}

	list, err := e.service.Set(c, req)
	if err != nil {
		zap.L().Error("get template list error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
	return
}

// List 设备电流配置列表
// @Summary 设备电流配置列表
// @Description 设备电流配置列表
// @Tags 系统管理/设备配置
// @Param  object query systemModels.SysDeviceElectricSettingDQL true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /system/electric/list [get]
func (e *Electric) List(c *gin.Context) {
	req := new(systemModels.SysDeviceElectricSettingDQL)
	err := c.ShouldBindQuery(req)
	if err != nil {
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

// Remove 删除设备电流配置
// @Summary 删除设备电流配置
// @Description 删除设备电流配置
// @Tags 系统管理/设备配置
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /system/electric/remove/{ids}  [delete]
func (e *Electric) Remove(c *gin.Context) {
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
