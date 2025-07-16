package alertController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/alert/alertDao"
)

type AlertTemplate struct {
	dao alertDao.AlertSinkTemplateDao
}

func NewAlertTemplate(dao alertDao.AlertSinkTemplateDao) *AlertTemplate {
	return &AlertTemplate{
		dao: dao,
	}
}

// List 网关协议列表
// @Summary 网关协议列表
// @Description 网关协议列表
// @Tags 网关管理/协议管理
// @Param  object query gatewayModels.SysSetGatewayInboundConfigReq true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /alert/template/list [get]
func (ac *AlertTemplate) List(c *gin.Context) {

}

// Set 网关协议列表
// @Summary 网关协议列表
// @Description 网关协议列表
// @Tags 网关管理/协议管理
// @Param  object query gatewayModels.SysSetGatewayInboundConfigReq true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /alert/template/set [get]
func (ac *AlertTemplate) Set(c *gin.Context) {

}

// Remove 网关协议列表
// @Summary 网关协议列表
// @Description 网关协议列表
// @Tags 网关管理/协议管理
// @Param  object query gatewayModels.SysSetGatewayInboundConfigReq true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /alert/template/remove [get]
func (ac *AlertTemplate) Remove(c *gin.Context) {

}
