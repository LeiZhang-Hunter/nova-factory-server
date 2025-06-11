package deviceController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type TemplateData struct {
	service         deviceService.ISysModbusDeviceConfigDataService
	templateService deviceService.IDeviceTemplateService
}

func NewTemplateData(service deviceService.ISysModbusDeviceConfigDataService, templateService deviceService.IDeviceTemplateService) *TemplateData {
	return &TemplateData{
		service:         service,
		templateService: templateService,
	}
}

func (t *TemplateData) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/asset/device/template/data")
	group.GET("/list", middlewares.HasPermission("asset:device:template:data:list"), t.List)
	group.POST("/set", middlewares.HasPermission("asset:device:template:data:set"), t.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("asset:device:template:data:remove"), t.Remove)
}

// List 获取模板数据列表
// @Summary 获取模板数据列表
// @Description 获取模板数据列表
// @Tags 资产管理/设备模板管理
// @Param  object query deviceModels.SysModbusDeviceConfigDataListReq true "获取模板数据列表请求参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/device/template/data/list [get]
func (t *TemplateData) List(c *gin.Context) {
	req := new(deviceModels.SysModbusDeviceConfigDataListReq)
	_ = c.ShouldBind(req)
	list, err := t.service.List(c, req)
	if err != nil {
		zap.L().Error("读取设备模板数据失败", zap.Error(err))
	}
	baizeContext.SuccessData(c, list)
}

// Set 设置设备模板板数信息
// @Summary 设置设备模板板数信息
// @Description 设置设备模板板数信息
// @Tags 资产管理/设备模板管理
// @Param  object body deviceModels.SetSysModbusDeviceConfigDataReq true "设置设备模板板数信息参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置模板成功"
// @Router /asset/device/template/data/set [post]
func (t *TemplateData) Set(c *gin.Context) {
	info := new(deviceModels.SetSysModbusDeviceConfigDataReq)
	err := c.ShouldBindJSON(info)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	tempInfo, err := t.templateService.GetById(c, info.TemplateID)
	if err != nil || tempInfo == nil {
		baizeContext.Waring(c, "模板不存在")
		return
	}
	if info.DeviceConfigID > 0 {
		vo, err := t.service.Update(c, info)
		if err != nil {
			zap.L().Error("修改设备模板数据失败", zap.Error(err))
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, vo)
	} else {
		vo, err := t.service.Add(c, info)
		if err != nil {
			zap.L().Error("插入设备模板数据失败", zap.Error(err))
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, vo)
	}
}

// Remove 删除设备数据模板
// @Summary 删除设备数据模板
// @Description 删除设备数据模板
// @Tags 资产管理/设备模板管理
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /asset/device/template/data/remove/{ids}  [delete]
func (t *TemplateData) Remove(c *gin.Context) {
	ids := baizeContext.ParamStringArray(c, "ids")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		return
	}
	err := t.service.Remove(c, ids)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}
	baizeContext.Success(c)
}
