// 管家婆全渠道控制器入口，汇总 OAuth、Token、API 三个子控制器。
// 通过 wire 注入服务接口，统一管理 /api/v1/erp-api/qqd 下的所有路由。
package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// ProviderSet 控制器的 wire 依赖集合，注入 Authorize/Token/API 并组装 Controller
var ProviderSet = wire.NewSet(NewAuthorize, NewToken, NewAPI, wire.Struct(new(Controller), "*"))

// Controller 管家婆全渠道顶层控制器，聚合授权、Token 签发和 API 调用三个模块
type Controller struct {
	Authorize *Authorize
	Token     *Token
	API       *API
}

// PublicRoutes 注册公开路由：/oauth/authorize、/oauth/token、/api
// POST /api 通过 method 参数分发到商品列表、商品新增、库存更新等子功能
func (q *Controller) PublicRoutes(group *gin.RouterGroup) {
	qqd := group.Group("/qqd")
	{
		qqd.GET("/oauth/authorize", q.Authorize.Authorize)
		qqd.POST("/oauth/token", q.Token.Token)
		qqd.POST("/api", q.API.API)
	}
}
