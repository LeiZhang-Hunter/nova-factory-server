package qcIqcController

import (
	"nova-factory-server/app/business/qcIqc/qcIqcApi"
	"nova-factory-server/app/business/qcIqc/qcIqcService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// QcIqcController 来料检验单控制器
type QcIqcController struct {
	qcIqcService qcIqcService.IQcIqcService
}

// NewQcIqcController 创建来料检验单控制器
func NewQcIqcController(qcIqcService qcIqcService.IQcIqcService) *QcIqcController {
	return &QcIqcController{
		qcIqcService: qcIqcService,
	}
}

// PrivateRoutes 私有路由
func (c *QcIqcController) PrivateRoutes(router *gin.RouterGroup) {
	qcIqc := router.Group("/qc/iqc")
	qcIqc.GET("/list", middlewares.HasPermission("qc:iqc:list"), c.GetQcIqcList)            // 获取来料检验单列表
	qcIqc.GET("/get/:iqcId", middlewares.HasPermission("qc:iqc:get"), c.GetQcIqcById)       // 根据ID获取来料检验单
	qcIqc.POST("/create", middlewares.HasPermission("qc:iqc:create"), c.CreateQcIqc)        // 创建来料检验单
	qcIqc.PUT("/update", middlewares.HasPermission("qc:iqc:update"), c.UpdateQcIqc)         // 更新来料检验单
	qcIqc.DELETE("/delete", middlewares.HasPermission("qc:iqc:delete"), c.DeleteQcIqcByIds) // 删除来料检验单
}

// GetQcIqcList 获取来料检验单列表
// @Summary 获取来料检验单列表
// @Description 分页查询来料检验单列表
// @Tags 来料检验单
// @Accept json
// @Produce json
// @Param pageNum query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param iqcCode query string false "检验单编号"
// @Param iqcName query string false "检验单名称"
// @Param templateId query int64 false "检验模板ID"
// @Param vendorId query int64 false "供应商ID"
// @Param vendorCode query string false "供应商编码"
// @Param vendorName query string false "供应商名称"
// @Param itemId query int64 false "产品物料ID"
// @Param itemCode query string false "产品物料编码"
// @Param itemName query string false "产品物料名称"
// @Param status query string false "单据状态"
// @Param inspector query string false "检测人员"
// @Router /qc/iqc/list [get]
func (c *QcIqcController) GetQcIqcList(ctx *gin.Context) {
	var req qcIqcApi.QcIqcQueryReq

	// 绑定查询参数
	if err := ctx.ShouldBindQuery(&req); err != nil {
		zap.L().Error("绑定查询参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	// 调用服务层
	result, err := c.qcIqcService.List(ctx, &req)
	if err != nil {
		zap.L().Error("查询来料检验单列表失败", zap.Error(err))
		baizeContext.Waring(ctx, "查询失败")
		return
	}

	// 返回结果
	baizeContext.SuccessData(ctx, result)
}

// GetQcIqcById 根据ID获取来料检验单
// @Summary 根据ID获取来料检验单
// @Description 根据ID获取来料检验单详情
// @Tags 来料检验单
// @Accept json
// @Produce json
// @Param iqcId path int true "来料检验单ID"
// @Router /qc/iqc/get/{iqcId} [get]
func (c *QcIqcController) GetQcIqcById(ctx *gin.Context) {
	iqcIdStr := ctx.Param("iqcId")
	iqcId, err := strconv.ParseInt(iqcIdStr, 10, 64)
	if err != nil {
		zap.L().Error("解析来料检验单ID失败", zap.String("iqcId", iqcIdStr), zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	iqc, err := c.qcIqcService.GetById(ctx, iqcId)
	if err != nil {
		zap.L().Error("查询来料检验单失败", zap.Int64("iqcId", iqcId), zap.Error(err))
		baizeContext.Waring(ctx, "查询失败")
		return
	}

	if iqc == nil {
		baizeContext.Waring(ctx, "来料检验单不存在")
		return
	}

	// 返回结果
	baizeContext.SuccessData(ctx, iqc)
}

// CreateQcIqc 创建来料检验单
// @Summary 创建来料检验单
// @Description 创建新的来料检验单
// @Tags 来料检验单
// @Accept json
// @Produce json
// @Param iqc body qcIqcApi.QcIqcCreateReq true "来料检验单信息"
// @Router /qc/iqc/create [post]
func (c *QcIqcController) CreateQcIqc(ctx *gin.Context) {
	var req qcIqcApi.QcIqcCreateReq

	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.L().Error("绑定创建参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	// 调用服务层
	_, err := c.qcIqcService.Create(ctx, &req)
	if err != nil {
		zap.L().Error("创建来料检验单失败", zap.Error(err))
		baizeContext.Waring(ctx, "创建失败")
		return
	}

	// 返回结果
	baizeContext.Success(ctx)
}

// UpdateQcIqc 更新来料检验单
// @Summary 更新来料检验单
// @Description 更新来料检验单信息
// @Tags 来料检验单
// @Accept json
// @Produce json
// @Param iqc body qcIqcApi.QcIqcUpdateReq true "来料检验单信息"
// @Router /qc/iqc/update [put]
func (c *QcIqcController) UpdateQcIqc(ctx *gin.Context) {
	var req qcIqcApi.QcIqcUpdateReq

	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.L().Error("绑定更新参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	// 调用服务层
	_, err := c.qcIqcService.Update(ctx, &req)
	if err != nil {
		zap.L().Error("更新来料检验单失败", zap.Error(err))
		baizeContext.Waring(ctx, "更新失败")
		return
	}

	// 返回结果
	baizeContext.Success(ctx)
}

// DeleteQcIqcByIds 批量删除来料检验单
// @Summary 批量删除来料检验单
// @Description 批量删除来料检验单
// @Tags 来料检验单
// @Accept json
// @Produce json
// @Param iqcIds body []int64 true "来料检验单ID列表"
// @Router /qc/iqc/delete [delete]
func (c *QcIqcController) DeleteQcIqcByIds(ctx *gin.Context) {
	var iqcIds []int64

	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&iqcIds); err != nil {
		zap.L().Error("绑定删除参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	if len(iqcIds) == 0 {
		baizeContext.Waring(ctx, "请选择要删除的来料检验单")
		return
	}

	// 调用服务层
	err := c.qcIqcService.Delete(ctx, iqcIds)
	if err != nil {
		zap.L().Error("删除来料检验单失败", zap.Error(err))
		baizeContext.Waring(ctx, "删除失败")
		return
	}

	// 返回结果
	baizeContext.Success(ctx)
}
