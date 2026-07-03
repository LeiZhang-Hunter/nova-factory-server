package shopcontroller

import (
	"nova-factory-server/app/business/shop/config/models"
	"nova-factory-server/app/business/shop/config/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

const (
	enterpriseKeyAccountName   = "enterprise_account_name"
	enterpriseKeyAccountNumber = "enterprise_account_number"
	enterpriseKeyAccountBank   = "enterprise_account_bank"
)

var enterpriseAllKeys = []string{
	enterpriseKeyAccountName, enterpriseKeyAccountNumber, enterpriseKeyAccountBank,
}

// EnterpriseAccountConfig 企业账户配置控制器
type EnterpriseAccountConfig struct {
	service service.IShopSysConfigService
}

// NewEnterpriseAccountConfig 创建企业账户配置控制器。
func NewEnterpriseAccountConfig(service service.IShopSysConfigService) *EnterpriseAccountConfig {
	return &EnterpriseAccountConfig{service: service}
}

// PrivateRoutes 注册企业账户配置路由
func (s *EnterpriseAccountConfig) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/config/enterprise-account")
	group.GET("/get", middlewares.HasPermission("shop:config:enterpriseAccount:query"), s.Get)
	group.PUT("/update", middlewares.HasPermission("shop:config:enterpriseAccount:edit"), s.Update)
}

// Get 获取企业账户配置
func (s *EnterpriseAccountConfig) Get(c *gin.Context) {
	data, err := s.service.GetByConfigKeys(c, enterpriseAllKeys)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, map[string]interface{}{"configs": data})
}

// Update 更新企业账户配置
func (s *EnterpriseAccountConfig) Update(c *gin.Context) {
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
