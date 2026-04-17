package shopcontroller

import (
	"nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/business/shop/user/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Address 商城用户地址控制器
type Address struct {
	service service.IShopAddressService
}

// NewAddress 创建商城用户地址控制器。
func NewAddress(service service.IShopAddressService) *Address {
	return &Address{service: service}
}

// PrivateRoutes 注册商城用户地址路由
func (s *Address) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/user/address")
	group.GET("/list", middlewares.HasPermission("shop:user:address:list"), s.List)
	group.POST("/set", middlewares.HasPermission("shop:user:address:set"), s.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("shop:user:address:remove"), s.Remove)
}

// List 获取商城用户地址列表
// @Summary 获取商城用户地址列表
// @Description 获取商城用户地址列表
// @Tags 商城/用户地址
// @Param object query models.AddressQuery true "商城用户地址查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/user/address/list [get]
func (s *Address) List(c *gin.Context) {
	req := new(models.AddressQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set 新增或修改商城用户地址
// @Summary 新增或修改商城用户地址
// @Description 新增或修改商城用户地址
// @Tags 商城/用户地址
// @Param object body models.AddressSetReq true "商城用户地址设置参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /shop/user/address/set [post]
func (s *Address) Set(c *gin.Context) {
	req := new(models.AddressSetReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.Set(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Remove 删除商城用户地址
// @Summary 删除商城用户地址
// @Description 根据ID删除商城用户地址
// @Tags 商城/用户地址
// @Param ids path string true "商城用户地址ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/user/address/remove/{ids} [delete]
func (s *Address) Remove(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := s.service.Remove(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
