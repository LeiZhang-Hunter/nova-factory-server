package mastercontroller

import (
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Supplier 供应商控制器
type Supplier struct {
	service masterservice.ISupplierService
}

// NewSupplier 创建供应商控制器。
func NewSupplier(service masterservice.ISupplierService) *Supplier {
	return &Supplier{service: service}
}

// PrivateRoutes 注册供应商私有路由。
func (o *Supplier) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/master/supplier")
	group.GET("/list", middlewares.HasPermission("erp:master:supplier:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:master:supplier:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:master:supplier:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:master:supplier:remove"), o.Delete)
}

// List 查询供应商列表。
// @Summary 查询供应商列表
// @Description 按条件分页查询供应商列表
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param object query mastermodels.SupplierQuery true "ERP 供应商查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/master/supplier/list [get]
func (o *Supplier) List(c *gin.Context) {
	req := new(mastermodels.SupplierQuery)
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

// GetByID 查询供应商详情。
// @Summary 查询供应商详情
// @Description 根据ID查询供应商详情
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/master/supplier/query/{id} [get]
func (o *Supplier) GetByID(c *gin.Context) {
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

// Set 新增或修改供应商。
// @Summary 新增或修改供应商
// @Description 新增或修改供应商
// @Tags ERP/基础资料
// @Security BearerAuth
// @Accept application/json
// @Param body body mastermodels.SupplierUpsert true "ERP 供应商参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/master/supplier/set [post]
func (o *Supplier) Set(c *gin.Context) {
	req := new(mastermodels.SupplierUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *mastermodels.Supplier
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

// Delete 删除供应商。
// @Summary 删除供应商
// @Description 根据ID删除供应商，多个ID用逗号分隔
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/master/supplier/remove/{ids} [delete]
func (o *Supplier) Delete(c *gin.Context) {
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
