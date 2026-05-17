package mastercontroller

import (
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Customer 客户控制器
type Customer struct {
	service masterservice.ICustomerService
}

// NewCustomer 创建客户控制器。
func NewCustomer(service masterservice.ICustomerService) *Customer {
	return &Customer{service: service}
}

// PrivateRoutes 注册客户私有路由。
func (o *Customer) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/master/customer")
	group.GET("/list", middlewares.HasPermission("erp:master:customer:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:master:customer:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:master:customer:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:master:customer:remove"), o.Delete)
}

// List 查询客户列表。
// @Summary 查询客户列表
// @Description 按条件分页查询客户列表
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param object query mastermodels.CustomerQuery true "ERP 客户查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/master/customer/list [get]
func (o *Customer) List(c *gin.Context) {
	req := new(mastermodels.CustomerQuery)
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

// GetByID 查询客户详情。
// @Summary 查询客户详情
// @Description 根据ID查询客户详情
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/master/customer/query/{id} [get]
func (o *Customer) GetByID(c *gin.Context) {
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

// Set 新增或修改客户。
// @Summary 新增或修改客户
// @Description 新增或修改客户
// @Tags ERP/基础资料
// @Security BearerAuth
// @Accept application/json
// @Param body body mastermodels.CustomerUpsert true "ERP 客户参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/master/customer/set [post]
func (o *Customer) Set(c *gin.Context) {
	req := new(mastermodels.CustomerUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *mastermodels.Customer
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

// Delete 删除客户。
// @Summary 删除客户
// @Description 根据ID删除客户，多个ID用逗号分隔
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/master/customer/remove/{ids} [delete]
func (o *Customer) Delete(c *gin.Context) {
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
