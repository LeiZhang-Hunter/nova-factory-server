package shopcontroller

import (
	"nova-factory-server/app/business/shop/config/models"
	"nova-factory-server/app/business/shop/config/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LogisticsConfig 物流配置控制器
type LogisticsConfig struct {
	service service.IShopLogisticsConfigService
}

// NewLogisticsConfig 创建物流配置控制器
func NewLogisticsConfig(service service.IShopLogisticsConfigService) *LogisticsConfig {
	return &LogisticsConfig{service: service}
}

// PrivateRoutes 注册物流配置私有路由
func (l *LogisticsConfig) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/config/logistics")
	group.GET("/list", middlewares.HasPermission("shop:config:logistics:list"), l.List)
	group.POST("/set", middlewares.HasPermission("shop:config:logistics:set"), l.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("shop:config:logistics:remove"), l.Remove)
}

func (l *LogisticsConfig) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/shop/config/logistics/list", "shop:config:logistics:list")
	router.RegisterPermission("POST", "/shop/config/logistics/set", "shop:config:logistics:set")
	router.RegisterPermission("DELETE", "/shop/config/logistics/remove/:ids", "shop:config:logistics:remove")
}

// List 查询物流配置列表
// @Summary 查询物流配置列表
// @Description 按条件分页查询物流配置
// @Tags 商城/物流配置
// @Security BearerAuth
// @Param object query models.ShopLogisticsConfigQuery true "物流配置查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /shop/config/logistics/list [get]
func (l *LogisticsConfig) List(c *gin.Context) {
	req := new(models.ShopLogisticsConfigQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := l.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set 新增或修改物流配置
// @Summary 新增或修改物流配置
// @Description 新增或修改物流配置
// @Tags 商城/物流配置
// @Security BearerAuth
// @Accept application/json
// @Param body body models.ShopLogisticsConfigSet true "物流配置参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /shop/config/logistics/set [post]
func (l *LogisticsConfig) Set(c *gin.Context) {
	req := new(models.ShopLogisticsConfigSet)
	if err := c.ShouldBindJSON(req); err != nil {
		zap.L().Error("param error", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *models.ShopLogisticsConfig
		err  error
	)
	if req.ID > 0 {
		data, err = l.service.Update(c, req)
	} else {
		data, err = l.service.Create(c, req)
	}
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Remove 删除物流配置
// @Summary 删除物流配置
// @Description 根据ID删除物流配置，多个ID用逗号分隔
// @Tags 商城/物流配置
// @Security BearerAuth
// @Param ids path string true "物流配置ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/config/logistics/{ids} [delete]
func (l *LogisticsConfig) Remove(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := l.service.DeleteByIDs(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
