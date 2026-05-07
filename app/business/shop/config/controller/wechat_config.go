package shopcontroller

import (
	"nova-factory-server/app/business/shop/config/models"
	"nova-factory-server/app/business/shop/config/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// WechatConfig 微信小程序配置控制器
type WechatConfig struct {
	service service.IShopSysConfigService
}

// NewWechatConfig 创建微信小程序配置控制器。
func NewWechatConfig(service service.IShopSysConfigService) *WechatConfig {
	return &WechatConfig{service: service}
}

// PrivateRoutes 注册微信小程序配置路由
func (s *WechatConfig) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/config/wechat")
	group.GET("/get", middlewares.HasPermission("shop:config:wechat:query"), s.GetWechatConfig)
	group.PUT("/update", middlewares.HasPermission("shop:config:wechat:edit"), s.UpdateWechatConfig)
}

// GetWechatConfig 获取微信小程序配置
// @Summary 获取微信小程序配置
// @Description 获取微信小程序配置
// @Tags 商城/系统配置
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/config/wechat/get [get]
func (s *WechatConfig) GetWechatConfig(c *gin.Context) {
	data, err := s.service.GetWechatConfig(c)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// UpdateWechatConfig 更新微信小程序配置
// @Summary 更新微信小程序配置
// @Description 更新微信小程序配置
// @Tags 商城/系统配置
// @Security BearerAuth
// @Produce application/json
// @Param object body models.WechatConfigReq true "微信小程序配置参数"
// @Success 200 {object} response.ResponseData "更新成功"
// @Router /shop/config/wechat/update [put]
func (s *WechatConfig) UpdateWechatConfig(c *gin.Context) {
	req := new(models.WechatConfigReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if err := s.service.UpdateWechatConfig(c, req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, "更新成功")
}
