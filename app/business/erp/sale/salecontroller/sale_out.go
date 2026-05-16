package salecontroller

import (
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/business/erp/sale/saleservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// SaleOut ERP 销售出库控制器
type SaleOut struct {
	service saleservice.ISaleOutService
}

// NewSaleOut 创建 ERP 销售出库控制器。
func NewSaleOut(service saleservice.ISaleOutService) *SaleOut {
	return &SaleOut{service: service}
}

// PrivateRoutes 注册 ERP 销售出库私有路由。
func (o *SaleOut) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/sale/out")
	group.GET("/list", middlewares.HasPermission("erp:sale:out:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:sale:out:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:sale:out:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:sale:out:remove"), o.Delete)
}

// List 查询 ERP 销售出库列表。
// @Summary 查询 ERP 销售出库列表
// @Description 按条件分页查询 ERP 销售出库列表
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param object query salemodels.SaleOutQuery true "ERP 销售出库查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/sale/out/list [get]
func (o *SaleOut) List(c *gin.Context) {
	req := new(salemodels.SaleOutQuery)
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

// GetByID 查询 ERP 销售出库详情。
// @Summary 查询 ERP 销售出库详情
// @Description 根据ID查询 ERP 销售出库详情
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/sale/out/query/{id} [get]
func (o *SaleOut) GetByID(c *gin.Context) {
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

// Set 新增或修改 ERP 销售出库。
// @Summary 新增或修改 ERP 销售出库
// @Description 新增或修改 ERP 销售出库
// @Tags ERP/销售管理
// @Security BearerAuth
// @Accept application/json
// @Param body body salemodels.SaleOutUpsert true "ERP 销售出库参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/sale/out/set [post]
func (o *SaleOut) Set(c *gin.Context) {
	req := new(salemodels.SaleOutUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *salemodels.SaleOut
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

// Delete 删除 ERP 销售出库。
// @Summary 删除 ERP 销售出库
// @Description 根据ID删除 ERP 销售出库，多个ID用逗号分隔
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/sale/out/remove/{ids} [delete]
func (o *SaleOut) Delete(c *gin.Context) {
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
