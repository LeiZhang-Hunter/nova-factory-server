package controller

import (
	"nova-factory-server/app/business/shop/logistics/models"
	"nova-factory-server/app/business/shop/logistics/service"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Tracking 物流轨迹查询控制器
type Tracking struct {
	service service.ITrackingService
}

// NewTracking 创建物流轨迹查询控制器
func NewTracking(service service.ITrackingService) *Tracking {
	return &Tracking{service: service}
}

// PrivateRoutes 注册管理端路由
func (t *Tracking) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/logistics/tracking")
	group.POST("/query", t.Query)
}

// AppRoutes 注册小程序端路由
func (t *Tracking) AppRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/logistics/tracking")
	group.POST("/query", t.Query)
}

// Query 即时查询物流轨迹
// @Summary 即时查询物流轨迹
// @Description 根据运单号和物流公司编码查询物流轨迹，优先从缓存获取
// @Tags 商城/物流查询
// @Security BearerAuth
// @Accept application/json
// @Param body body models.TrackingQueryRequest true "物流轨迹查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /shop/logistics/tracking/query [post]
func (t *Tracking) Query(c *gin.Context) {
	req := new(models.TrackingQueryRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		zap.L().Error("param error", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}

	resp, err := t.service.Query(c, req.Outsid, req.CompanyCode)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, resp)
}
