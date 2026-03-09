package cameraController

import (
	"fmt"
	"net/http"
	"nova-factory-server/app/business/iot/asset/camera/cameraModels"
	"nova-factory-server/app/business/iot/asset/camera/cameraService"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CameraController 摄像头控制器
type CameraController struct {
	cameraService cameraService.ICameraService
}

// NewCameraController 创建摄像头控制器实例
func NewCameraController(cameraService cameraService.ICameraService) *CameraController {
	return &CameraController{cameraService: cameraService}
}

// PrivateRoutes 注册路由
func (c *CameraController) PrivateRoutes(router *gin.RouterGroup) {
	cameraGroup := router.Group("/asset/camera")
	{
		cameraGroup.POST("/set", c.Set)
		cameraGroup.DELETE("/remove/:id", c.Delete)
		cameraGroup.GET("/info/:id", c.GetById)
		cameraGroup.GET("/list", c.List)
	}
}

// Set 设置摄像头信息
func (c *CameraController) Set(ctx *gin.Context) {
	var camera cameraModels.IotCamera
	if err := ctx.ShouldBindJSON(&camera); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if camera.Id > 0 {
		if err := c.cameraService.Update(&camera); err != nil {
			zap.L().Error("cameraService create error", zap.Error(err))
			baizeContext.Waring(ctx, "设置摄像头失败")
			return
		}
		baizeContext.SuccessData(ctx, camera)
		return
	}

	if err := c.cameraService.Create(&camera); err != nil {
		zap.L().Error("cameraService create error", zap.Error(err))
		baizeContext.Waring(ctx, "设置摄像头失败")
		return
	}
	baizeContext.SuccessData(ctx, camera)
}

// Delete 删除摄像头
func (c *CameraController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID不能为空"})
		return
	}

	if err := c.cameraService.Delete(ParseInt64(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetById 根据ID获取摄像头
func (c *CameraController) GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID不能为空"})
		return
	}

	camera, err := c.cameraService.GetById(ParseInt64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, camera)
}

// List 获取摄像头列表
func (c *CameraController) List(ctx *gin.Context) {
	cameras, err := c.cameraService.List()
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}

	baizeContext.SuccessData(ctx, cameras)
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
