package qcOqcController

import (
	"strconv"

	"nova-factory-server/app/business/qcOqc/qcOqcApi"
	"nova-factory-server/app/business/qcOqc/qcOqcService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// QcOqcController 出货检验单控制器
type QcOqcController struct {
	qcOqcService qcOqcService.IQcOqcService
}

// NewQcOqcController 创建出货检验单控制器
func NewQcOqcController(qcOqcService qcOqcService.IQcOqcService) *QcOqcController {
	return &QcOqcController{
		qcOqcService: qcOqcService,
	}
}

// PrivateRoutes 私有路由
func (c *QcOqcController) PrivateRoutes(router *gin.RouterGroup) {
	qcOqc := router.Group("/qc/oqc")
	qcOqc.GET("/list", middlewares.HasPermission("qc:oqc:list"), c.GetQcOqcList)       // 获取出货检验单列表
	qcOqc.GET("/get/:oqcId", middlewares.HasPermission("qc:oqc:get"), c.GetQcOqcById)  // 根据ID获取出货检验单
	qcOqc.POST("/create", middlewares.HasPermission("qc:oqc:create"), c.CreateQcOqc)   // 创建出货检验单
	qcOqc.PUT("/update", middlewares.HasPermission("qc:oqc:update"), c.UpdateQcOqc)    // 更新出货检验单
	qcOqc.DELETE("/delete", middlewares.HasPermission("qc:oqc:delete"), c.DeleteQcOqc) // 删除出货检验单
}

// GetQcOqcList 获取出货检验单列表
// @Summary 获取出货检验单列表
// @Description 分页查询出货检验单列表
// @Tags 出货检验单
// @Accept json
// @Produce json
// @Param pageNum query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param oqcCode query string false "检验单编号"
// @Param oqcName query string false "检验单名称"
// @Param templateId query int64 false "检验模板ID"
// @Param sourceDocId query int64 false "来源单据ID"
// @Param sourceDocType query string false "来源单据类型"
// @Param sourceDocCode query string false "来源单据编号"
// @Param sourceLineId query int64 false "来源单据行ID"
// @Param clientId query int64 false "客户ID"
// @Param clientCode query string false "客户编码"
// @Param clientName query string false "客户名称"
// @Param batchCode query string false "批次号"
// @Param itemId query int64 false "产品物料ID"
// @Param itemCode query string false "产品物料编码"
// @Param itemName query string false "产品物料名称"
// @Param specification query string false "规格型号"
// @Param unitOfMeasure query string false "单位"
// @Param quantityMinCheck query float64 false "最低检测数"
// @Param quantityMaxUnqualified query float64 false "最大不合格数"
// @Param quantityOut query float64 false "发货数量"
// @Param quantityCheck query float64 false "本次检测数量"
// @Param quantityUnqualified query float64 false "不合格数"
// @Param quantityQualified query float64 false "合格数量"
// @Param crRate query float64 false "致命缺陷率"
// @Param majRate query float64 false "严重缺陷率"
// @Param minRate query float64 false "轻微缺陷率"
// @Param crQuantity query float64 false "致命缺陷数量"
// @Param majQuantity query float64 false "严重缺陷数量"
// @Param minQuantity query float64 false "轻微缺陷数量"
// @Param checkResult query string false "检测结果"
// @Param outDate query string false "出货日期"
// @Param inspectDate query string false "检测日期"
// @Param inspector query string false "检测人员"
// @Param status query string false "单据状态"
// @Param createBy query string false "创建者"
// @Router /qc/oqc/list [get]
func (c *QcOqcController) GetQcOqcList(ctx *gin.Context) {
	var req qcOqcApi.QcOqcQueryReq

	// 绑定查询参数
	if err := ctx.ShouldBindQuery(&req); err != nil {
		zap.L().Error("绑定查询参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	// 调用服务查询数据
	list, total, err := c.qcOqcService.GetQcOqcList(ctx, &req)
	if err != nil {
		zap.L().Error("查询出货检验单列表失败", zap.Error(err))
		baizeContext.Waring(ctx, "查询失败")
		return
	}

	// 返回结果
	baizeContext.SuccessData(ctx, gin.H{
		"rows":  list,
		"total": total,
	})
}

// GetQcOqcById 根据ID获取出货检验单
// @Summary 根据ID获取出货检验单
// @Description 根据ID获取出货检验单详情
// @Tags 出货检验单
// @Accept json
// @Produce json
// @Param oqcId path int true "出货检验单ID"
// @Router /qc/oqc/get/{oqcId} [get]
func (c *QcOqcController) GetQcOqcById(ctx *gin.Context) {
	// 获取ID参数
	oqcIdStr := ctx.Param("oqcId")
	oqcId, err := strconv.ParseInt(oqcIdStr, 10, 64)
	if err != nil {
		zap.L().Error("解析出货检验单ID失败", zap.String("oqcId", oqcIdStr), zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	// 调用服务查询数据
	oqc, err := c.qcOqcService.GetQcOqcById(ctx, oqcId)
	if err != nil {
		zap.L().Error("查询出货检验单失败", zap.Int64("oqcId", oqcId), zap.Error(err))
		baizeContext.Waring(ctx, "查询失败")
		return
	}

	if oqc == nil {
		baizeContext.Waring(ctx, "出货检验单不存在")
		return
	}

	// 返回结果
	baizeContext.SuccessData(ctx, oqc)
}

// CreateQcOqc 创建出货检验单
// @Summary 创建出货检验单
// @Description 创建新的出货检验单
// @Tags 出货检验单
// @Accept json
// @Produce json
// @Param oqc body qcOqcApi.QcOqcCreateReq true "出货检验单信息"
// @Router /qc/oqc/create [post]
func (c *QcOqcController) CreateQcOqc(ctx *gin.Context) {
	var req qcOqcApi.QcOqcCreateReq

	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.L().Error("绑定创建参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	// 调用服务创建数据
	err := c.qcOqcService.CreateQcOqc(ctx, &req)
	if err != nil {
		zap.L().Error("创建出货检验单失败", zap.Error(err))
		baizeContext.Waring(ctx, "创建失败")
		return
	}

	// 返回结果
	baizeContext.Success(ctx)
}

// UpdateQcOqc 更新出货检验单
// @Summary 更新出货检验单
// @Description 更新出货检验单信息
// @Tags 出货检验单
// @Accept json
// @Produce json
// @Param oqc body qcOqcApi.QcOqcUpdateReq true "出货检验单信息"
// @Router /qc/oqc/update [put]
func (c *QcOqcController) UpdateQcOqc(ctx *gin.Context) {
	var req qcOqcApi.QcOqcUpdateReq

	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.L().Error("绑定更新参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	// 调用服务更新数据
	err := c.qcOqcService.UpdateQcOqc(ctx, &req)
	if err != nil {
		zap.L().Error("更新出货检验单失败", zap.Error(err))
		baizeContext.Waring(ctx, "更新失败")
		return
	}

	// 返回结果
	baizeContext.Success(ctx)
}

// DeleteQcOqc 删除出货检验单
// @Summary 删除出货检验单
// @Description 批量删除出货检验单
// @Tags 出货检验单
// @Accept json
// @Produce json
// @Param oqcIds body []int64 true "出货检验单ID列表"
// @Router /qc/oqc/delete [delete]
func (c *QcOqcController) DeleteQcOqc(ctx *gin.Context) {
	var req qcOqcApi.QcOqcDeleteReq

	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.L().Error("绑定删除参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	// 调用服务删除数据
	err := c.qcOqcService.DeleteQcOqc(ctx, &req)
	if err != nil {
		zap.L().Error("删除出货检验单失败", zap.Error(err))
		baizeContext.Waring(ctx, "删除失败")
		return
	}

	// 返回结果
	baizeContext.Success(ctx)
}
