package modules

import "github.com/gin-gonic/gin"

type Controller interface {
	// PrivateRoutes 注册需要登陆的路由
	PrivateRoutes(router *gin.RouterGroup)

	// PublicRoutes 注册不需要注册的路由
	PublicRoutes(router *gin.RouterGroup)
}
