package mastercontroller

import (
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

// ProductCategory ERP 产品分类控制器
type ProductCategory struct {
	service masterservice.IProductCategoryService
}

// NewProductCategory 创建 ERP 产品分类控制器。
func NewProductCategory(service masterservice.IProductCategoryService) *ProductCategory {
	return &ProductCategory{service: service}
}

// PrivateRoutes 注册 ERP 产品分类私有路由。
func (o *ProductCategory) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/master/product-category")
	group.GET("/list", middlewares.HasPermission("erp:master:productCategory:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:master:productCategory:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:master:productCategory:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:master:productCategory:remove"), o.Delete)
}

func (o *ProductCategory) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/erp/master/product-category/list", "erp:master:productCategory:list")
	router.RegisterPermission("GET", "/erp/master/product-category/query/:id", "erp:master:productCategory:query")
	router.RegisterPermission("POST", "/erp/master/product-category/set", "erp:master:productCategory:set")
	router.RegisterPermission("DELETE", "/erp/master/product-category/remove/:ids", "erp:master:productCategory:remove")
}

// List 查询 ERP 产品分类列表。
// @Summary 查询 ERP 产品分类列表
// @Description 按条件分页查询 ERP 产品分类列表
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param object query mastermodels.ProductCategoryQuery true "ERP 产品分类查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/master/product-category/list [get]
func (o *ProductCategory) List(c *gin.Context) {
	req := new(mastermodels.ProductCategoryQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetByID 查询 ERP 产品分类详情。
// @Summary 查询 ERP 产品分类详情
// @Description 根据ID查询 ERP 产品分类详情
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/master/product-category/query/{id} [get]
func (o *ProductCategory) GetByID(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.GetByID(c, id)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set 新增或修改 ERP 产品分类。
// @Summary 新增或修改 ERP 产品分类
// @Description 新增或修改 ERP 产品分类
// @Tags ERP/基础资料
// @Security BearerAuth
// @Accept application/json
// @Param body body mastermodels.ProductCategoryUpsert true "ERP 产品分类参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/master/product-category/set [post]
func (o *ProductCategory) Set(c *gin.Context) {
	req := new(mastermodels.ProductCategoryUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *mastermodels.ProductCategory
		err  error
	)
	if req.ID > 0 {
		data, err = o.service.Update(c, req)
	} else {
		data, err = o.service.Create(c, req)
	}
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Delete 删除 ERP 产品分类。
// @Summary 删除 ERP 产品分类
// @Description 根据ID删除 ERP 产品分类，多个ID用逗号分隔
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/master/product-category/remove/{ids} [delete]
func (o *ProductCategory) Delete(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := o.service.DeleteByIDs(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
