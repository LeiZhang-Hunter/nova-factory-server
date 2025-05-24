package gatewayController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/gateway/gatewayModels"
	"nova-factory-server/app/business/gateway/gatewayService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Protocol struct {
	service gatewayService.ISysGatewayInboundConfigService
}

func NewProtocol(service gatewayService.ISysGatewayInboundConfigService) *Protocol {
	return &Protocol{
		service: service,
	}
}

func (p *Protocol) PrivateRoutes(router *gin.RouterGroup) {
	routers := router.Group("/gateway")
	routers.GET("/protocol/list", middlewares.HasPermission("gateway:protocol"), p.GetProtocolList)                 // 网关协议列表
	routers.POST("/protocol/set", middlewares.HasPermission("gateway:protocol:set"), p.SetProtocol)                 // 设置网关协议
	routers.DELETE("/protocol/remove/:ids", middlewares.HasPermission("gateway:protocol:remove"), p.RemoveProtocol) //移除网关协议
}

// GetProtocolList 网关协议列表
// @Summary 网关协议列表
// @Description 网关协议列表
// @Tags 网关管理/协议管理
// @Param  object query gatewayModels.SysSetGatewayInboundConfigReq true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /gateway/protocol/list [get]
func (p *Protocol) GetProtocolList(c *gin.Context) {
	req := new(gatewayModels.SysSetGatewayInboundConfigReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := p.service.List(c, req)
	if err != nil {
		return
	}
	baizeContext.SuccessData(c, list)
	return
}

// SetProtocol 设置网关协议
// @Summary 设置网关协议
// @Description 设置网关协议
// @Tags 网关管理/协议管理
// @Param  object query gatewayModels.SysSetGatewayInboundConfig true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /gateway/protocol/set [post]
func (p *Protocol) SetProtocol(c *gin.Context) {
	data := new(gatewayModels.SysSetGatewayInboundConfig)
	err := c.ShouldBindJSON(data)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if data.GatewayConfigID == 0 {
		ret, err := p.service.Add(c, data)
		if err != nil {
			return
		}
		baizeContext.SuccessData(c, ret)
	} else {
		ret, err := p.service.Update(c, data)
		if err != nil {
			return
		}
		baizeContext.SuccessData(c, ret)
	}

	return
}

// RemoveProtocol 删除网关协议
// @Summary 删除网关协议
// @Description 删除网关协议
// @Tags 网关管理/协议管理
// @Param  ids path string true "ids"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /gateway/protocol/remove [delete]
func (p *Protocol) RemoveProtocol(c *gin.Context) {
	ids := baizeContext.ParamStringArray(c, "ids")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		return
	}
	err := p.service.Remove(c, ids)
	if err != nil {
		zap.L().Error("删除失败", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
	return
}
