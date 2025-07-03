package qcRqcController

import (
	"strconv"

	"nova-factory-server/app/business/qcRqc/qcRqcApi"
	"nova-factory-server/app/business/qcRqc/qcRqcService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// QcRqcController 退料检验单控制器
type QcRqcController struct {
	qcRqcService qcRqcService.IQcRqcService
}

// NewQcRqcController 创建退料检验单控制器
func NewQcRqcController(qcRqcService qcRqcService.IQcRqcService) *QcRqcController {
	return &QcRqcController{qcRqcService: qcRqcService}
}

// PrivateRoutes 私有路由
func (c *QcRqcController) PrivateRoutes(router *gin.RouterGroup) {
	qcRqc := router.Group("/qc/rqc")
	qcRqc.GET("/list", middlewares.HasPermission("qc:rqc:list"), c.GetQcRqcList)            // 获取退料检验单列表
	qcRqc.GET("/get/:rqcId", middlewares.HasPermission("qc:rqc:get"), c.GetQcRqcById)       // 根据ID获取退料检验单
	qcRqc.POST("/create", middlewares.HasPermission("qc:rqc:create"), c.CreateQcRqc)        // 创建退料检验单
	qcRqc.PUT("/update", middlewares.HasPermission("qc:rqc:update"), c.UpdateQcRqc)         // 更新退料检验单
	qcRqc.DELETE("/delete", middlewares.HasPermission("qc:rqc:delete"), c.DeleteQcRqcByIds) // 删除退料检验单
}

// GetQcRqcList 获取退料检验单列表
// @Summary 获取退料检验单列表
// @Description 分页查询退料检验单列表
// @Tags 退料检验单
// @Accept json
// @Produce json
// @Param pageNum query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param rqcCode query string false "检验单编号"
// @Param rqcName query string false "检验单名称"
// @Param templateId query int64 false "检验模板ID"
// @Param sourceDocId query int64 false "来源单据ID"
// @Param sourceDocType query string false "来源单据类型"
// @Param sourceDocCode query string false "来源单据编号"
// @Param sourceLineId query int64 false "来源单据行ID"
// @Param itemId query int64 false "产品物料ID"
// @Param itemCode query string false "产品物料编码"
// @Param itemName query string false "产品物料名称"
// @Param specification query string false "规格型号"
// @Param unitOfMeasure query string false "单位"
// @Param unitName query string false "单位名称"
// @Param batchId query int64 false "批次ID"
// @Param batchCode query string false "批次号"
// @Param quantityCheck query float64 false "检测数量"
// @Param quantityUnqualified query float64 false "不合格数"
// @Param quantityQualified query float64 false "合格品数量"
// @Param checkResult query string false "检测结果"
// @Param inspectDate query string false "检测日期"
// @Param userId query int64 false "检测人员ID"
// @Param userName query string false "检测人员名称"
// @Param nickName query string false "检测人员"
// @Param status query string false "单据状态"
// @Param createBy query string false "创建者"
// @Router /qc/rqc/list [get]
func (c *QcRqcController) GetQcRqcList(ctx *gin.Context) {
	var req qcRqcApi.QcRqcQueryReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		zap.L().Error("绑定查询参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	// 设置默认分页参数
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	list, total, err := c.qcRqcService.GetQcRqcList(ctx, &req)
	if err != nil {
		zap.L().Error("查询退料检验单列表失败", zap.Error(err))
		baizeContext.Waring(ctx, "查询失败")
		return
	}

	// 返回结果
	baizeContext.SuccessData(ctx, gin.H{
		"rows":  list,
		"total": total,
	})
}

// GetQcRqcById 根据ID获取退料检验单
// @Summary 根据ID获取退料检验单
// @Description 根据ID获取退料检验单详情
// @Tags 退料检验单
// @Accept json
// @Produce json
// @Param rqcId path int true "退料检验单ID"
// @Router /qc/rqc/get/{rqcId} [get]
func (c *QcRqcController) GetQcRqcById(ctx *gin.Context) {
	rqcIdStr := ctx.Param("rqcId")
	rqcId, err := strconv.ParseInt(rqcIdStr, 10, 64)
	if err != nil {
		baizeContext.ParameterError(ctx)
		return
	}

	rqc, err := c.qcRqcService.GetQcRqcById(ctx, rqcId)
	if err != nil {
		zap.L().Error("查询退料检验单失败", zap.Error(err))
		baizeContext.Waring(ctx, "查询失败")
		return
	}

	if rqc == nil {
		baizeContext.Waring(ctx, "退料检验单不存在")
		return
	}

	// 返回结果
	baizeContext.SuccessData(ctx, rqc)
}

// CreateQcRqc 创建退料检验单
// @Summary 创建退料检验单
// @Description 创建新的退料检验单
// @Tags 退料检验单
// @Accept json
// @Produce json
// @Param rqc body qcRqcApi.QcRqcCreateReq true "退料检验单信息"
// @Router /qc/rqc/create [post]
func (c *QcRqcController) CreateQcRqc(ctx *gin.Context) {
	var req qcRqcApi.QcRqcCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.L().Error("绑定创建参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	err := c.qcRqcService.CreateQcRqc(ctx, &req)
	if err != nil {
		zap.L().Error("创建退料检验单失败", zap.Error(err))
		baizeContext.Waring(ctx, "创建失败")
		return
	}

	// 返回结果
	baizeContext.Success(ctx)
}

// UpdateQcRqc 更新退料检验单
// @Summary 更新退料检验单
// @Description 更新退料检验单信息
// @Tags 退料检验单
// @Accept json
// @Produce json
// @Param rqc body qcRqcApi.QcRqcUpdateReq true "退料检验单信息"
// @Router /qc/rqc/update [put]
func (c *QcRqcController) UpdateQcRqc(ctx *gin.Context) {
	var req qcRqcApi.QcRqcUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.L().Error("绑定更新参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	err := c.qcRqcService.UpdateQcRqc(ctx, &req)
	if err != nil {
		zap.L().Error("更新退料检验单失败", zap.Error(err))
		baizeContext.Waring(ctx, "更新失败")
		return
	}

	// 返回结果
	baizeContext.Success(ctx)
}

// DeleteQcRqcByIds 批量删除退料检验单
// @Summary 批量删除退料检验单
// @Description 批量删除退料检验单
// @Tags 退料检验单
// @Accept json
// @Produce json
// @Param rqcIds body []int64 true "退料检验单ID列表"
// @Router /qc/rqc/delete [delete]
func (c *QcRqcController) DeleteQcRqcByIds(ctx *gin.Context) {
	var rqcIds []int64
	if err := ctx.ShouldBindJSON(&rqcIds); err != nil {
		zap.L().Error("绑定删除参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	err := c.qcRqcService.DeleteQcRqcByIds(ctx, rqcIds)
	if err != nil {
		zap.L().Error("删除退料检验单失败", zap.Error(err))
		baizeContext.Waring(ctx, "删除失败")
		return
	}

	// 返回结果
	baizeContext.Success(ctx)
}
