package cameraController

import (
	"fmt"
	"nova-factory-server/app/business/iot/asset/camera/cameraModels"
	"nova-factory-server/app/business/iot/asset/camera/cameraService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Camera 摄像头控制器
type Camera struct {
	cameraService cameraService.ICameraService
}

// NewCameraController 创建摄像头控制器实例
func NewCameraController(cameraService cameraService.ICameraService) *Camera {
	return &Camera{
		cameraService: cameraService,
	}
}

// PrivateRoutes 注册路由
func (c *Camera) PrivateRoutes(router *gin.RouterGroup) {
	cameraGroup := router.Group("/iot/asset/camera")
	cameraGroup.POST("/set", middlewares.HasPermission("iot:asset:camera:set"), c.Set)
	cameraGroup.DELETE("/remove/:ids", middlewares.HasPermission("iot:asset:camera:remove"), c.Delete)
	cameraGroup.GET("/info/:id", middlewares.HasPermission("iot:asset:camera:info"), c.GetById)
	cameraGroup.GET("/list", middlewares.HasPermission("iot:asset:camera:list"), c.List)
}

// Set 设置摄像头信息
func (c *Camera) Set(ctx *gin.Context) {
	var camera cameraModels.IotCamera
	if err := ctx.ShouldBindJSON(&camera); err != nil {
		zap.L().Error("参数错误:", zap.Error(err))
		baizeContext.Waring(ctx, "参数错误")
		return
	}

	if camera.Id > 0 {
		if err := c.cameraService.Update(ctx, &camera); err != nil {
			zap.L().Error("cameraService create error", zap.Error(err))
			baizeContext.Waring(ctx, "设置摄像头失败")
			return
		}
		baizeContext.SuccessData(ctx, camera)
		return
	}

	if err := c.cameraService.Create(ctx, &camera); err != nil {
		zap.L().Error("cameraService create error", zap.Error(err))
		baizeContext.Waring(ctx, "设置摄像头失败")
		return
	}
	baizeContext.SuccessData(ctx, camera)
}

// Delete 删除摄像头
func (c *Camera) Delete(ctx *gin.Context) {
	ids := baizeContext.ParamStringArray(ctx, "ids")
	if len(ids) == 0 {
		baizeContext.Waring(ctx, "ID不能为空")
		return
	}

	if err := c.cameraService.Delete(ids); err != nil {
		zap.L().Error("cameraService delete error", zap.Error(err))
		baizeContext.Waring(ctx, err.Error())
		return
	}

	baizeContext.Success(ctx)
}

// GetById 根据ID获取摄像头
func (c *Camera) GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		baizeContext.Waring(ctx, "ID不能为空")
		return
	}

	camera, err := c.cameraService.GetById(ParseInt64(id))
	if err != nil {
		zap.L().Error("cameraService get error", zap.Error(err))
		baizeContext.Waring(ctx, err.Error())
		return
	}

	baizeContext.SuccessData(ctx, camera)
}

// List 获取摄像头列表
// @Summary 读取摄像头列表
// @Description 读取摄像头列表
// @Tags 资产管理/摄像头管理
// @Param  object query cameraModels.IotCameraListReq true "摄像头列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/camera/list [get]
func (c *Camera) List(ctx *gin.Context) {
	req := new(cameraModels.IotCameraListReq)
	err := ctx.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(ctx)
		return
	}

	cameraList, err := c.cameraService.List(req)
	if err != nil {
		zap.L().Error("get camera list error", zap.Error(err))
		baizeContext.Waring(ctx, err.Error())
		return
	}

	baizeContext.SuccessData(ctx, cameraList)
}

// ParseInt64 字符串转int64
func ParseInt64(s string) int64 {
	var i int64
	_, err := fmt.Sscanf(s, "%d", &i)
	if err != nil {
		return 0
	}
	return i
}
