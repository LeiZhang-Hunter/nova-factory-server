package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// ProviderSet WMS 仓储控制器 Provider。
var ProviderSet = wire.NewSet(NewWarehouseArea, NewWarehouseLocation, wire.Struct(new(Controller), "*"))

// Controller WMS 仓储控制器聚合。
type Controller struct {
	WarehouseArea     *WarehouseArea
	WarehouseLocation *WarehouseLocation
}

// PrivateRoutes 注册 WMS 仓储私有路由。
func (c *Controller) PrivateRoutes(router *gin.RouterGroup) {
	if c.WarehouseArea != nil {
		c.WarehouseArea.PrivateRoutes(router)
	}
	if c.WarehouseLocation != nil {
		c.WarehouseLocation.PrivateRoutes(router)
	}
}
