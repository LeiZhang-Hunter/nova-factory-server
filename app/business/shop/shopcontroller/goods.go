package shopcontroller

import (
	"go.uber.org/zap"
	"nova-factory-server/app/business/shop/shopmodels"
	"nova-factory-server/app/business/shop/shopservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Goods struct {
	service shopservice.IShopGoodsService
}

func NewGoods(service shopservice.IShopGoodsService) *Goods {
	return &Goods{service: service}
}

// PrivateRoutes 注册商品路由
func (s *Goods) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/goods")
	group.GET("/list", middlewares.HasPermission("shop:goods:list"), s.List)
	group.GET("/info/:id", middlewares.HasPermission("shop:goods:info"), s.GetByID)
	group.POST("/add", middlewares.HasPermission("shop:goods:add"), s.Create)
	group.PUT("/edit", middlewares.HasPermission("shop:goods:edit"), s.Update)
	group.DELETE("/remove/:ids", middlewares.HasPermission("shop:goods:remove"), s.Delete)
}

// PublicRoutes 导入接口注册
func (s *Goods) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/goods")
	group.Any("/import", s.Import)
}

// List 获取商品列表
// @Summary 获取商品列表
// @Description 获取商品列表
// @Tags 商城/商品管理
// @Param object query shopmodels.GoodsQuery true "商品查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/goods/list [get]
func (s *Goods) List(c *gin.Context) {
	req := new(shopmodels.GoodsQuery)
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

// GetByID 获取商品详情
// @Summary 获取商品详情
// @Description 根据ID获取商品详情
// @Tags 商城/商品管理
// @Param id path int true "商品ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/goods/{id} [get]
func (s *Goods) GetByID(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.GetByID(c, id)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Create 新增商品
// @Summary 新增商品
// @Description 新增商品
// @Tags 商城/商品管理
// @Param object body shopmodels.GoodsUpsert true "商品新增参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "新增成功"
// @Router /shop/goods [post]
func (s *Goods) Create(c *gin.Context) {
	req := new(shopmodels.GoodsUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.Create(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Update 修改商品
// @Summary 修改商品
// @Description 修改商品
// @Tags 商城/商品管理
// @Param object body shopmodels.GoodsUpsert true "商品修改参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "修改成功"
// @Router /shop/goods [put]
func (s *Goods) Update(c *gin.Context) {
	req := new(shopmodels.GoodsUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.ID == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.Update(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Delete 删除商品
// @Summary 删除商品
// @Description 根据ID删除商品
// @Tags 商城/商品管理
// @Param ids path string true "商品ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/goods/{ids} [delete]
func (s *Goods) Delete(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := s.service.DeleteByIDs(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// Import 导入商品
// @Summary 导入商品
// @Description 导入商品
// @Tags 商城/商品管理
// @Param object body shopmodels.GoodsUpsert true "商品新增参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "新增成功"
// @Router /shop/goods/import [post]
func (s *Goods) Import(c *gin.Context) {
	goods := new(shopmodels.ImportGoodsList)
	err := c.ShouldBindJSON(goods)
	if err != nil {
		zap.L().Error("parse goods fail", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	if err = s.service.Import(c, goods.Records); err != nil {
		zap.L().Error("import goods fail", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
