package shopcontroller

import (
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Sku struct {
	service shopservice.IShopSkuService
}

func NewSku(service shopservice.IShopSkuService) *Sku {
	return &Sku{service: service}
}

// PrivateRoutes 注册商品规格路由
func (s *Sku) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/sku")
	group.GET("/list", middlewares.HasPermission("shop:sku:list"), s.List)
	group.GET("/info/:id", middlewares.HasPermission("shop:sku:query"), s.GetByID)
	group.POST("/add", middlewares.HasPermission("shop:sku:add"), s.Create)
	group.PUT("/update", middlewares.HasPermission("shop:sku:update"), s.Update)
	group.DELETE("/remove/:ids", middlewares.HasPermission("shop:sku:remove"), s.Delete)
}

// List 获取商品规格列表
// @Summary 获取商品规格列表
// @Description 获取商品规格列表
// @Tags 商城/商品规格
// @Param object query shopmodels.GoodsSkuQuery true "商品规格查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/sku/list [get]
func (s *Sku) List(c *gin.Context) {
	req := new(shopmodels.GoodsSkuQuery)
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

// GetByID 获取商品规格详情
// @Summary 获取商品规格详情
// @Description 根据ID获取商品规格详情
// @Tags 商城/商品规格
// @Param id path int true "商品规格ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/sku/{id} [get]
func (s *Sku) GetByID(c *gin.Context) {
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

// Create 新增商品规格
// @Summary 新增商品规格
// @Description 新增商品规格
// @Tags 商城/商品规格
// @Param object body shopmodels.GoodsSkuUpsert true "商品规格新增参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "新增成功"
// @Router /shop/sku [post]
func (s *Sku) Create(c *gin.Context) {
	req := new(shopmodels.GoodsSkuUpsert)
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

// Update 修改商品规格
// @Summary 修改商品规格
// @Description 修改商品规格
// @Tags 商城/商品规格
// @Param object body shopmodels.GoodsSkuUpsert true "商品规格修改参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "修改成功"
// @Router /shop/sku [put]
func (s *Sku) Update(c *gin.Context) {
	req := new(shopmodels.GoodsSkuUpsert)
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

// Delete 删除商品规格
// @Summary 删除商品规格
// @Description 根据ID删除商品规格
// @Tags 商城/商品规格
// @Param ids path string true "商品规格ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/sku/{ids} [delete]
func (s *Sku) Delete(c *gin.Context) {
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
