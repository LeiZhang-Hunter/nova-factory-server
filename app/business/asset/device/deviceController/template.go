package deviceController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/constant/device"
	"nova-factory-server/app/constant/protocols"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Template struct {
	service deviceService.IDeviceTemplateService
}

func NewTemplate(d deviceService.IDeviceTemplateService) (*Template, error) {
	return &Template{
		service: d,
	}, nil
}

func (t *Template) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/asset/device/template")
	group.GET("/list", middlewares.HasPermission("asset:device:template:list"), t.List)
	group.POST("/set", middlewares.HasPermission("asset:device:template:set"), t.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("asset:device:template:remove"), t.Remove)
	group.GET("/protocols", middlewares.HasPermission("asset:device:template:protocols"), t.Protocols)
}

// List 获取模板列表
// @Summary 获取模板列表
// @Description 获取模板列表
// @Tags 资产管理/设备模板管理
// @Param  object query deviceModels.SysDeviceTemplateDQL true "获取模板列表请求参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/device/template/list [get]
func (t *Template) List(c *gin.Context) {
	req := new(deviceModels.SysDeviceTemplateDQL)
	_ = c.ShouldBind(req)
	list, err := t.service.List(c, req)
	if err != nil {
		zap.L().Error("读取设备模板失败", zap.Error(err))
	}
	baizeContext.SuccessData(c, list)
}

// Protocols 读取协议列表
// @Summary 读取协议列表
// @Description 读取协议列表
// @Tags 资产管理/设备协议管理
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/device/template/protocols [get]
func (t *Template) Protocols(c *gin.Context) {
	baizeContext.SuccessData(c, []string{
		protocols.MODBUS,
	})
}

// Set 设置设备模板信息
// @Summary 设置设备模板信息
// @Description 设置设备模板信息
// @Tags 资产管理/设备模板管理
// @Param  object body deviceModels.SysDeviceTemplateSetReq true "设备模板参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置模板成功"
// @Router /asset/device/template/set [post]
func (t *Template) Set(c *gin.Context) {
	info := new(deviceModels.SysDeviceTemplateSetReq)
	err := c.ShouldBindJSON(info)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	// 校验协议类型是否合法
	if !device.CheckProtocol(info.Protocol) {
		baizeContext.Waring(c, "协议类型错误")
		return
	}

	if info.TemplateID > 0 {
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

// Remove 删除设备模板
// @Summary 删除设备模板
// @Description 删除设备模板
// @Tags 资产管理/设备模板管理
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /asset/device/template/remove/{ids}  [delete]
func (t *Template) Remove(c *gin.Context) {
	ids := baizeContext.ParamStringArray(c, "ids")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		return
	}
	err := t.service.Remove(c, ids)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
