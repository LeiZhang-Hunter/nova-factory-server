package settingcontroller

import (
	"nova-factory-server/app/business/erp/setting/settingmodels"
	"nova-factory-server/app/business/erp/setting/settingservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logistics ERP物流公司控制器
type Logistics struct {
	service settingservice.ILogisticsCompanyService
}

// NewLogistics 创建 ERP物流公司控制器。
func NewLogistics(service settingservice.ILogisticsCompanyService) *Logistics {
	return &Logistics{service: service}
}

// PrivateRoutes 注册 ERP物流公司私有路由。
func (l *Logistics) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/setting/logistics-company")
	group.GET("/list", middlewares.HasPermission("erp:setting:logisticsCompany:list"), l.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:setting:logisticsCompany:query"), l.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:setting:logisticsCompany:set"), l.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:setting:logisticsCompany:remove"), l.Delete)
}

// List 查询 ERP物流公司列表。
// @Summary 查询 ERP物流公司列表
// @Description 按条件分页查询 ERP物流公司列表
// @Tags ERP/系统配置
// @Security BearerAuth
// @Param object query settingmodels.LogisticsCompanyQuery true "ERP物流公司查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/setting/logistics-company/list [get]
func (l *Logistics) List(c *gin.Context) {
	req := new(settingmodels.LogisticsCompanyQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := l.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetByID 查询 ERP物流公司详情。
// @Summary 查询 ERP物流公司详情
// @Description 根据ID查询 ERP物流公司详情
// @Tags ERP/系统配置
// @Security BearerAuth
// @Param id path int true "物流公司ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/setting/logistics-company/{id} [get]
func (l *Logistics) GetByID(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := l.service.GetByID(c, id)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set 新增或修改 ERP物流公司。
// @Summary 新增或修改 ERP物流公司
// @Description 新增或修改 ERP物流公司
// @Tags ERP/系统配置
// @Security BearerAuth
// @Accept application/json
// @Param body body settingmodels.LogisticsCompanyUpsert true "ERP物流公司参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/setting/logistics-company/set [post]
func (l *Logistics) Set(c *gin.Context) {
	req := new(settingmodels.LogisticsCompanyUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		zap.L().Error("param error", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *settingmodels.LogisticsCompany
		err  error
	)
	if req.ID > 0 {
		data, err = l.service.Update(c, req)
	} else {
		data, err = l.service.Create(c, req)
	}
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Delete 删除 ERP物流公司。
// @Summary 删除 ERP物流公司
// @Description 根据ID删除 ERP物流公司，多个ID用逗号分隔
// @Tags ERP/系统配置
// @Security BearerAuth
// @Param ids path string true "物流公司ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/setting/logistics-company/{ids} [delete]
func (l *Logistics) Delete(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := l.service.DeleteByIDs(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
