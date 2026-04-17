package shopcontroller

import (
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Address 商城地址控制器
type Address struct {
	service shopservice.IShopAddressService
}

// NewAddress 创建商城地址控制器。
func NewAddress(service shopservice.IShopAddressService) *Address {
	return &Address{service: service}
}

// PrivateRoutes 注册商城地址路由
func (s *Address) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/address")
	group.GET("/list", middlewares.HasPermission("shop:address:list"), s.List)
	group.POST("/set", middlewares.HasPermission("shop:address:set"), s.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("shop:address:remove"), s.Remove)
}

// List 获取商城地址列表
// @Summary 获取商城地址列表
// @Description 获取商城地址列表
// @Tags 商城/地址管理
// @Param object query shopmodels.AddressQuery true "商城地址查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/address/list [get]
func (s *Address) List(c *gin.Context) {
	req := new(shopmodels.AddressQuery)
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

// Set 新增或修改商城地址
// @Summary 新增或修改商城地址
// @Description 新增或修改商城地址
// @Tags 商城/地址管理
// @Param object body shopmodels.AddressSetReq true "商城地址设置参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /shop/address/set [post]
func (s *Address) Set(c *gin.Context) {
	req := new(shopmodels.AddressSetReq)
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

// Remove 删除商城地址
// @Summary 删除商城地址
// @Description 根据ID删除商城地址
// @Tags 商城/地址管理
// @Param ids path string true "商城地址ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/address/remove/{ids} [delete]
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
