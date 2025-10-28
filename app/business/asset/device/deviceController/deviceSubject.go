package deviceController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/middlewares"
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
	subject := router.Group("/asset/deviceGroup")
	subject.GET("/list", middlewares.HasPermission("asset:deviceGroup"), ds.List)                  // 设备列表
	subject.POST("/set", middlewares.HasPermission("asset:deviceGroup:set"), ds.Set)               // 设置设备信息
	subject.DELETE("/:groupIds", middlewares.HasPermission("asset:deviceGroup:remove"), ds.Remove) //删除设备分组列表
}

func (ds *DeviceSubject) List(c *gin.Context) {

}

func (ds *DeviceSubject) Set(c *gin.Context) {

}

func (ds *DeviceSubject) Remove(c *gin.Context) {

}
