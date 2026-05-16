package mastercontroller

import (
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Account ERP 结算账户控制器
type Account struct {
	service masterservice.IAccountService
}

// NewAccount 创建 ERP 结算账户控制器。
func NewAccount(service masterservice.IAccountService) *Account {
	return &Account{service: service}
}

// PrivateRoutes 注册 ERP 结算账户私有路由。
func (o *Account) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/master/account")
	group.GET("/list", middlewares.HasPermission("erp:master:account:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:master:account:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:master:account:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:master:account:remove"), o.Delete)
}

// List 查询 ERP 结算账户列表。
// @Summary 查询 ERP 结算账户列表
// @Description 按条件分页查询 ERP 结算账户列表
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param object query mastermodels.AccountQuery true "ERP 结算账户查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/master/account/list [get]
func (o *Account) List(c *gin.Context) {
	req := new(mastermodels.AccountQuery)
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

// GetByID 查询 ERP 结算账户详情。
// @Summary 查询 ERP 结算账户详情
// @Description 根据ID查询 ERP 结算账户详情
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/master/account/query/{id} [get]
func (o *Account) GetByID(c *gin.Context) {
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

// Set 新增或修改 ERP 结算账户。
// @Summary 新增或修改 ERP 结算账户
// @Description 新增或修改 ERP 结算账户
// @Tags ERP/基础资料
// @Security BearerAuth
// @Accept application/json
// @Param body body mastermodels.AccountUpsert true "ERP 结算账户参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/master/account/set [post]
func (o *Account) Set(c *gin.Context) {
	req := new(mastermodels.AccountUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *mastermodels.Account
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

// Delete 删除 ERP 结算账户。
// @Summary 删除 ERP 结算账户
// @Description 根据ID删除 ERP 结算账户，多个ID用逗号分隔
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/master/account/remove/{ids} [delete]
func (o *Account) Delete(c *gin.Context) {
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
