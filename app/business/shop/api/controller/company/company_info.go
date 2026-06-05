package company

import (
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// CompanyInfo provides public mini app company information.
type CompanyInfo struct {
	service service.IApiShopCompanyInfoService
}

// NewCompanyInfo creates a company info controller.
func NewCompanyInfo(service service.IApiShopCompanyInfoService) *CompanyInfo {
	return &CompanyInfo{service: service}
}

// PublicRoutes registers public company info routes.
func (cc *CompanyInfo) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/company")
	group.GET("/info", cc.Info)
}

// Info gets public company information.
// @Summary 获取小程序公司信息
// @Description 获取小程序关于我们页面展示的公司信息
// @Tags app接口/商城/App公司信息
// @Produce application/json
// @Success 200 {object} response.ResponseData{data=models.CompanyInfoResp} "获取成功"
// @Router /api/v1/app/shop/company/info [get]
func (cc *CompanyInfo) Info(c *gin.Context) {
	data, err := cc.service.Get(c)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

var _ = models.CompanyInfoResp{}
