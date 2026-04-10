package shopcontroller

import (
	"nova-factory-server/app/business/shop/shopmodels"
	"nova-factory-server/app/business/shop/shopservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Category struct {
	service shopservice.IShopCategoryService
}

func NewCategory(service shopservice.IShopCategoryService) *Category {
	return &Category{service: service}
}

// PrivateRoutes 注册商品分类路由
func (s *Category) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/category")
	group.GET("/list", middlewares.HasPermission("shop:category:list"), s.List)
	group.GET("/:id", middlewares.HasPermission("shop:category:query"), s.GetByID)
	group.POST("", middlewares.HasPermission("shop:category:add"), s.Create)
	group.PUT("", middlewares.HasPermission("shop:category:edit"), s.Update)
	group.DELETE("/:ids", middlewares.HasPermission("shop:category:remove"), s.Delete)
}

// List 获取商品分类列表
// @Summary 获取商品分类列表
// @Description 获取商品分类列表
// @Tags 商城/商品分类
// @Param object query shopmodels.CategoryQuery true "商品分类查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/category/list [get]
func (s *Category) List(c *gin.Context) {
	req := new(shopmodels.CategoryQuery)
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

// GetByID 获取商品分类详情
// @Summary 获取商品分类详情
// @Description 根据ID获取商品分类详情
// @Tags 商城/商品分类
// @Param id path int true "商品分类ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/category/{id} [get]
func (s *Category) GetByID(c *gin.Context) {
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

// Create 新增商品分类
// @Summary 新增商品分类
// @Description 新增商品分类
// @Tags 商城/商品分类
// @Param object body shopmodels.CategoryUpsert true "商品分类新增参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "新增成功"
// @Router /shop/category [post]
func (s *Category) Create(c *gin.Context) {
	req := new(shopmodels.CategoryUpsert)
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

// Update 修改商品分类
// @Summary 修改商品分类
// @Description 修改商品分类
// @Tags 商城/商品分类
// @Param object body shopmodels.CategoryUpsert true "商品分类修改参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "修改成功"
// @Router /shop/category [put]
func (s *Category) Update(c *gin.Context) {
	req := new(shopmodels.CategoryUpsert)
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

// Delete 删除商品分类
// @Summary 删除商品分类
// @Description 根据ID删除商品分类
// @Tags 商城/商品分类
// @Param ids path string true "商品分类ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/category/{ids} [delete]
func (s *Category) Delete(c *gin.Context) {
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
