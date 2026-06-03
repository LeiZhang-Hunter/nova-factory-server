package shopcontroller

import (
	"nova-factory-server/app/business/shop/config/models"
	"nova-factory-server/app/business/shop/config/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

const (
	wechatKeyAppID                = "wechat_mini_program_app_id"
	wechatKeyAppSecret            = "wechat_mini_program_app_secret"
	wechatKeyToken                = "wechat_mini_program_token"
	wechatKeyEncodingAESKey       = "wechat_mini_program_encoding_aes_key"
	wechatKeyPayMchID             = "wechat_pay_mch_id"
	wechatKeyPayNotifyURL         = "wechat_pay_notify_url"
	wechatKeyPayAPIv3Key          = "wechat_pay_api_v3_key"
	wechatKeyPaySerialNo          = "wechat_pay_serial_no"
	wechatKeyPayPrivateKey        = "wechat_pay_private_key_path"
	wechatKeyPayPlatformPublicKey = "wechat_pay_platform_public_key_path"
)

var wechatAllKeys = []string{
	wechatKeyAppID, wechatKeyAppSecret, wechatKeyToken, wechatKeyEncodingAESKey,
	wechatKeyPayMchID, wechatKeyPayNotifyURL,
	wechatKeyPayAPIv3Key, wechatKeyPaySerialNo, wechatKeyPayPrivateKey, wechatKeyPayPlatformPublicKey,
}

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
	data, err := s.service.GetByConfigKeys(c, wechatAllKeys)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, map[string]interface{}{"configs": data})
}

// UpdateWechatConfig 更新微信小程序配置
// @Summary 更新微信小程序配置
// @Description 更新微信小程序配置
// @Tags 商城/系统配置
// @Security BearerAuth
// @Produce application/json
// @Param object body models.BatchConfigReq true "配置参数"
// @Success 200 {object} response.ResponseData "更新成功"
// @Router /shop/config/wechat/update [put]
func (s *WechatConfig) UpdateWechatConfig(c *gin.Context) {
	req := new(models.BatchConfigReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if err := s.service.BatchUpdate(c, req.Configs); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, "更新成功")
}
