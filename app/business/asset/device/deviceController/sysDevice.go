package deviceController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"
)

type DeviceInfo struct {
	iDeviceService deviceService.IDeviceService
}

func NewDeviceInfo(d deviceService.IDeviceService) *DeviceInfo {
	return &DeviceInfo{
		iDeviceService: d,
	}
}

func (di *DeviceInfo) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterSchema("GET", "/asset/device/list", deviceModels.DeviceListReq{}, nil)
}

func (di *DeviceInfo) PrivateRoutes(router *gin.RouterGroup) {
	asset := router.Group("/asset/device")
	asset.GET("/list", middlewares.HasPermission("asset:device"), di.GetDeviceList)                // 设备列表
	asset.POST("/set", middlewares.HasPermission("asset:device:set"), di.SetDevice)                // 设置设备信息
	asset.DELETE("/:deviceIds", middlewares.HasPermission("asset:device:remove"), di.DeviceRemove) // 设置设备信息
}

func (di *DeviceInfo) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/device")
	// 根据标签读取设备数据
	group.GET("/metric/tag", di.GetDeviceMetricByTag)
}

// GetDeviceList 获取设备列表
// @Summary 获取设备列表
// @Description 获取设备列表
// @Tags 资产管理
// @Param  object query deviceModels.DeviceListReq true "设备列表请求参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/device/list [get]
func (di *DeviceInfo) GetDeviceList(c *gin.Context) {
	req := new(deviceModels.DeviceListReq)
	err := c.ShouldBindQuery(req)
	list, err := di.iDeviceService.SelectDeviceList(c, req)
	if err != nil {
		zap.L().Error("读取设备分组失败", zap.Error(err))
	}
	baizeContext.SuccessData(c, list)
}

// SetDevice 设置设备信息
// @Summary 设置设备信息
// @Description 设置设备信息
// @Tags 资产管理
// @Param  object body deviceModels.DeviceInfo true "设备参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置成功"
// @Router /asset/device/set [post]
func (di *DeviceInfo) SetDevice(c *gin.Context) {
	info := new(deviceModels.DeviceInfo)
	err := c.ShouldBindJSON(info)
	if err != nil {
		zap.L().Error("set device param error", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	if info.DeviceId > 0 {
		vo, err := di.iDeviceService.UpdateDevice(c, info)
		if err != nil {
			zap.L().Error("更新设备数据失败", zap.Error(err))
			baizeContext.Waring(c, "更新设备数据失败")
			return
		}
		baizeContext.SuccessData(c, vo)
	} else {
		vo, err := di.iDeviceService.InsertDevice(c, info)
		if err != nil {
			zap.L().Error("插入设备数据失败", zap.Error(err))
			baizeContext.Waring(c, "插入设备数据失败")
			return
		}
		baizeContext.SuccessData(c, vo)
	}
}

// DeviceRemove 删除设备
// @Summary 删除设备
// @Description 删除设备
// @Tags 资产管理
// @Param  deviceIds path string true "deviceIds"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /asset/device/{deviceIds}  [delete]
func (di *DeviceInfo) DeviceRemove(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "deviceIds")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		return
	}
	err := di.iDeviceService.DeleteByDeviceIds(c, ids)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}
	baizeContext.Success(c)

}

// GetDeviceMetricByTag 获取设备数据通过标签
// @Summary 获取设备数据通过标签
// @Description 获取设备数据通过标签
// @Tags 资产管理
// @Param  object query deviceModels.DeviceListReq true "设备列表请求参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /api/v1/device/metric/tag [get]
func (di *DeviceInfo) GetDeviceMetricByTag(c *gin.Context) {
	req := new(deviceModels.DeviceTagListReq)
	err := c.ShouldBindQuery(req)
	list, err := di.iDeviceService.GetMetricByTag(c, req)
	if err != nil {
		zap.L().Error("读取设备分组失败", zap.Error(err))
	}
	baizeContext.SuccessData(c, list)
}
