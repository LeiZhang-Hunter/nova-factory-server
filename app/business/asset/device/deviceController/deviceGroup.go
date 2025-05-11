package deviceController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type DeviceGroup struct {
	iDeviceGroupService deviceService.IDeviceGroupService
}

func NewDeviceGroup(d deviceService.IDeviceGroupService) *DeviceGroup {
	return &DeviceGroup{
		iDeviceGroupService: d,
	}
}

func (di *DeviceGroup) PrivateRoutes(router *gin.RouterGroup) {
	asset := router.Group("/asset/deviceGroup")
	asset.GET("/list", middlewares.HasPermission("asset:deviceGroup"), di.GetDeviceGroupList)               // 设备列表
	asset.POST("/set", middlewares.HasPermission("asset:deviceGroup:set"), di.SetDeviceGroup)               // 设置设备信息
	asset.DELETE("/:groupIds", middlewares.HasPermission("asset:deviceGroup:remove"), di.DeviceGroupRemove) //删除设备分组列表
}

// GetDeviceGroupList 获取设备分组列表
// @Summary 获取设备分组列表
// @Description 获取设备分组列表
// @Tags 资产管理
// @Param  object query deviceModels.DeviceGroupDQL true "设备分组列表请求参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/deviceGroup/list [get]
func (di *DeviceGroup) GetDeviceGroupList(c *gin.Context) {
	req := new(deviceModels.DeviceGroupDQL)
	_ = c.ShouldBind(req)
	list, err := di.iDeviceGroupService.SelectDeviceGroupList(c, req)
	if err != nil {
		zap.L().Error("读取设备分组失败", zap.Error(err))
	}
	baizeContext.SuccessData(c, list)
}

// SetDeviceGroup 设置设备分组信息
// @Summary 设置设备分组信息
// @Description 设置设备分组信息
// @Tags 资产管理
// @Param  object body deviceModels.DeviceGroup true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /asset/deviceGroup/set [post]
func (di *DeviceGroup) SetDeviceGroup(c *gin.Context) {
	info := new(deviceModels.DeviceGroup)
	err := c.ShouldBindJSON(info)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if info.Name == nil || *info.Name == "" {
		baizeContext.Waring(c, "设备分组名不能是空")
		return
	}
	if info.GroupId > 0 {
		vo, err := di.iDeviceGroupService.UpdateDeviceGroup(c, info)
		if err != nil {
			zap.L().Error("修改设备分组数据失败", zap.Error(err))
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, vo)
	} else {
		vo, err := di.iDeviceGroupService.InsertDeviceGroup(c, info)
		if err != nil {
			zap.L().Error("插入设备分组数据失败", zap.Error(err))
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, vo)
	}
}

// DeviceGroupRemove 删除设备组
// @Summary 删除设备组
// @Description 删除设备组
// @Tags 资产管理
// @Param  groupIds path string true "groupIds"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /asset/deviceGroup/{groupIds}  [delete]
func (di *DeviceGroup) DeviceGroupRemove(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "groupIds")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		return
	}
	err := di.iDeviceGroupService.DeleteByGroupIds(c, ids)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}
	baizeContext.Success(c)

}
