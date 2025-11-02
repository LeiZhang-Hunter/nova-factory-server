package deviceController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

// DeviceSubject 设备点检保养项目表
type DeviceSubject struct {
	service deviceService.IDeviceSubjectService
}

func NewDeviceSubject(service deviceService.IDeviceSubjectService) *DeviceSubject {
	return &DeviceSubject{
		service: service,
	}
}

func (ds *DeviceSubject) PrivateRoutes(router *gin.RouterGroup) {
	subject := router.Group("/asset/device/subject")
	subject.GET("/list", middlewares.HasPermission("asset:device:subject:list"), ds.List)               // 点检保养项目列表
	subject.POST("/set", middlewares.HasPermission("asset:device:subject:set"), ds.Set)                 // 设置点检保养项目
	subject.DELETE("/remove/:ids", middlewares.HasPermission("asset:device:subject:remove"), ds.Remove) //删除点检保养项目
}

// List 点检保养项目列表
// @Summary 点检保养项目列表
// @Description 点检保养项目列表
// @Tags 设备管理/点检保养项目列表
// @Param  object query deviceModels.SysDeviceSubjectReq true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/device/subject/list [get]
func (ds *DeviceSubject) List(c *gin.Context) {
	req := new(deviceModels.SysDeviceSubjectReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	list, err := ds.service.List(c, req)
	if err != nil {
		zap.L().Error("读取点检保养项目失败", zap.Error(err))
		baizeContext.Waring(c, "读取点检保养项目失败")
		return
	}
	baizeContext.SuccessData(c, list)
}

// Set 设置点检保养项目
// @Summary 设置点检保养项目
// @Description 设置点检保养项目
// @Tags 设备管理/设置点检保养项目
// @Param  object body deviceModels.SysDeviceSubjectVO true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/device/subject/set [post]
func (ds *DeviceSubject) Set(c *gin.Context) {
	vo := new(deviceModels.SysDeviceSubjectVO)
	err := c.ShouldBindJSON(vo)
	if err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	if vo.ID == 0 {
		info, err := ds.service.GetBySubjectCode(c, vo.SubjectCode)
		if err != nil {
			baizeContext.Waring(c, err.Error())
			return
		}
		if info != nil {
			baizeContext.Waring(c, "设置点检保养项目编码存在")
			return
		}
	} else {
		info, err := ds.service.GetBySubjectCodeByNotId(c, vo.ID, vo.SubjectCode)
		if err != nil {
			baizeContext.Waring(c, err.Error())
			return
		}
		if info != nil {
			baizeContext.Waring(c, "设置点检保养项目编码存在")
			return
		}
	}
	set, err := ds.service.Set(c, vo)
	if err != nil {
		zap.L().Error("设置点检保养项目失败", zap.Error(err))
		baizeContext.Waring(c, "设置点检保养项目失败")
		return
	}
	baizeContext.SuccessData(c, set)
}

// Remove 删除点检保养项目
// @Summary 删除点检保养项目
// @Description 删除点检保养项目
// @Tags 设备管理/删除点检保养项目
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /asset/device/subject/remove/{ids}  [delete]
func (ds *DeviceSubject) Remove(c *gin.Context) {
	ids := baizeContext.ParamStringArray(c, "ids")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		baizeContext.ParameterError(c)
		return
	}
	err := ds.service.Remove(c, ids)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}

	baizeContext.Success(c)
}
