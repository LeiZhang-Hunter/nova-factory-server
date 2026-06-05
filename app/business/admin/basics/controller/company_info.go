package controller

import (
	"nova-factory-server/app/business/admin/basics/models"
	"nova-factory-server/app/business/admin/basics/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"strings"

	"github.com/gin-gonic/gin"
)

type CompanyInfo struct {
	cs service.ICompanyInfoService
}

func NewCompanyInfo(cs service.ICompanyInfoService) *CompanyInfo {
	return &CompanyInfo{cs: cs}
}

func (cc *CompanyInfo) PrivateRoutes(router *gin.RouterGroup) {
	basicsCompany := router.Group("/basics/company")
	basicsCompany.GET("/info", middlewares.HasPermission("basics:company:query"), cc.CompanyInfoGet)
	basicsCompany.PUT("/info", middlewares.SetLog("公司信息", middlewares.Update), middlewares.HasPermission("basics:company:edit"), cc.CompanyInfoSave)
}

// CompanyInfoGet 获取公司信息
// @Summary 获取公司信息
// @Description 获取公司信息
// @Tags 公司信息
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData{data=models.CompanyInfoVo} "成功"
// @Router /basics/company/info [get]
func (cc *CompanyInfo) CompanyInfoGet(c *gin.Context) {
	baizeContext.SuccessData(c, cc.cs.SelectCompanyInfo(c))
}

// CompanyInfoSave 保存公司信息
// @Summary 保存公司信息
// @Description 保存公司信息
// @Tags 公司信息
// @Param object body models.CompanyInfoVo true "公司信息"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "成功"
// @Router /basics/company/info [put]
func (cc *CompanyInfo) CompanyInfoSave(c *gin.Context) {
	company := new(models.CompanyInfoVo)
	if err := c.ShouldBindJSON(company); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	company.CompanyName = strings.TrimSpace(company.CompanyName)
	if company.CompanyName == "" {
		baizeContext.ParameterError(c)
		return
	}
	cc.cs.SaveCompanyInfo(c, company)
	baizeContext.Success(c)
}
