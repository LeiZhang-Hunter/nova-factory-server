package deviceController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type DeviceCheckMachinery struct {
	service deviceService.IDeviceCheckMachineryService
}

func NewDeviceCheckMachinery(service deviceService.IDeviceCheckMachineryService) *DeviceCheckMachinery {
	return &DeviceCheckMachinery{
		service: service,
	}
}

func (dm *DeviceCheckMachinery) PrivateRoutes(router *gin.RouterGroup) {
	subject := router.Group("/asset/device/check/machinery")
	subject.GET("/list", middlewares.HasPermission("asset:device:check:machinery:list"), dm.List)               // 点检保养项目列表
	subject.POST("/set", middlewares.HasPermission("asset:device:check:machinery:set"), dm.Set)                 // 设置点检保养项目
	subject.DELETE("/remove/:ids", middlewares.HasPermission("asset:device:check:machinery:remove"), dm.Remove) //删除点检保养项目
}

// List 设备清单列表
// @Summary 设备清单列表
// @Description 设备清单列表
// @Tags 设备管理/点检保养计划/设备清单
// @Param  object query deviceModels.SysDeviceCheckMachineryReq true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/device/check/machinery/list [get]
func (dm *DeviceCheckMachinery) List(c *gin.Context) {
	req := new(deviceModels.SysDeviceCheckMachineryReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	list, err := dm.service.List(c, req)
	if err != nil {
		zap.L().Error("读取点检保养计划失败", zap.Error(err))
		baizeContext.Waring(c, "读取点检保养计划失败")
		return
	}
	baizeContext.SuccessData(c, list)
}

// Set 设置设备清单
// @Summary 设置设备清单
// @Description 设置设备清单
// @Tags 设备管理/点检保养计划/设备清单
// @Param  object body deviceModels.SysDeviceCheckPlanVO true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/device/check/machinery/set [post]
func (dm *DeviceCheckMachinery) Set(c *gin.Context) {
	vo := new(deviceModels.SysDeviceCheckMachineryVO)
	err := c.ShouldBindJSON(vo)
	if err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	set, err := dm.service.Set(c, vo)
	if err != nil {
		zap.L().Error("设置设备清单失败", zap.Error(err))
		baizeContext.Waring(c, "设置设备清单失败")
		return
	}
	baizeContext.SuccessData(c, set)
}

// Remove 删除设备清单
// @Summary 删除设备清单
// @Description 删除设备清单
// @Tags 设备管理/点检保养计划/设备清单
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /asset/device/check/machinery/remove/{ids}  [delete]
func (dm *DeviceCheckMachinery) Remove(c *gin.Context) {
	ids := baizeContext.ParamStringArray(c, "ids")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		baizeContext.ParameterError(c)
		return
	}
	err := dm.service.Remove(c, ids)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}

	baizeContext.Success(c)
}
