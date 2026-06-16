package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"nova-factory-server/app/business/shop/order/service"
	"nova-factory-server/app/utils/baizeContext"
)

// AutoCancel 订单超时自动取消兜底控制器。
//
// 提供 HTTP 接口供 Linux Cron 等外部定时任务调用，与 Consumer 共享 ProcessExpiredOrders 逻辑。
type AutoCancel struct {
	service service.IOrderTimeoutService
}

// NewAutoCancel 创建订单超时自动取消控制器。
func NewAutoCancel(service service.IOrderTimeoutService) *AutoCancel {
	return &AutoCancel{service: service}
}

// PublicRoutes 注册兜底取消路由到 publicGroup（不走业务鉴权，用 token 防护）。
func (a *AutoCancel) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/order")
	group.GET("/auto_cancel", a.AutoCancel)
}

// AutoCancel 扫描并取消所有到期未支付订单。
// @Summary 订单超时自动取消兜底接口
// @Description 供外部定时任务调用，扫描 Redis 延迟队列取消到期订单
// @Tags Shop/订单
// @Param token query string true "防护令牌"
// @Produce application/json
// @Success 200 {object} response.ResponseData "取消成功，data 为取消数量"
// @Router /api/v1/app/shop/order/auto_cancel [get]
func (a *AutoCancel) AutoCancel(c *gin.Context) {
	expected := strings.TrimSpace(viper.GetString("shop.order.auto_cancel_token"))
	token := strings.TrimSpace(c.Query("token"))
	if expected == "" || token != expected {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	count, err := a.service.ProcessExpiredOrders(c.Request.Context())
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, count)
}
