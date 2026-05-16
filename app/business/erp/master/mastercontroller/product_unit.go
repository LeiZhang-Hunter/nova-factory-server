package mastercontroller

import (
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// ProductUnit ERP 产品单位控制器
type ProductUnit struct {
	service masterservice.IProductUnitService
}

// NewProductUnit 创建 ERP 产品单位控制器。
func NewProductUnit(service masterservice.IProductUnitService) *ProductUnit {
	return &ProductUnit{service: service}
}

// PrivateRoutes 注册 ERP 产品单位私有路由。
func (o *ProductUnit) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/master/product-unit")
	group.GET("/list", middlewares.HasPermission("erp:master:productUnit:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:master:productUnit:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:master:productUnit:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:master:productUnit:remove"), o.Delete)
}

// List 查询 ERP 产品单位列表。
// @Summary 查询 ERP 产品单位列表
// @Description 按条件分页查询 ERP 产品单位列表
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param object query mastermodels.ProductUnitQuery true "ERP 产品单位查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/master/product-unit/list [get]
func (o *ProductUnit) List(c *gin.Context) {
	req := new(mastermodels.ProductUnitQuery)
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

// GetByID 查询 ERP 产品单位详情。
// @Summary 查询 ERP 产品单位详情
// @Description 根据ID查询 ERP 产品单位详情
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/master/product-unit/query/{id} [get]
func (o *ProductUnit) GetByID(c *gin.Context) {
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

// Set 新增或修改 ERP 产品单位。
// @Summary 新增或修改 ERP 产品单位
// @Description 新增或修改 ERP 产品单位
// @Tags ERP/基础资料
// @Security BearerAuth
// @Accept application/json
// @Param body body mastermodels.ProductUnitUpsert true "ERP 产品单位参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/master/product-unit/set [post]
func (o *ProductUnit) Set(c *gin.Context) {
	req := new(mastermodels.ProductUnitUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *mastermodels.ProductUnit
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

// Delete 删除 ERP 产品单位。
// @Summary 删除 ERP 产品单位
// @Description 根据ID删除 ERP 产品单位，多个ID用逗号分隔
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/master/product-unit/remove/{ids} [delete]
func (o *ProductUnit) Delete(c *gin.Context) {
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
