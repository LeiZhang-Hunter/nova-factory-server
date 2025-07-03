package batchController

import (
	"nova-factory-server/app/business/batch/batchApi"
	"nova-factory-server/app/business/batch/batchService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BatchController struct {
	batchService batchService.IBatchService
}

func NewBatchController(batchService batchService.IBatchService) *BatchController {
	return &BatchController{
		batchService: batchService,
	}
}

func (c *BatchController) PrivateRoutes(router *gin.RouterGroup) {
	batch := router.Group("/batch")
	batch.GET("/list", middlewares.HasPermission("batch:list"), c.GetBatchList)
	batch.GET("/get/:batchId", middlewares.HasPermission("batch:get"), c.GetBatchById)
	batch.POST("/create", middlewares.HasPermission("batch:create"), c.CreateBatch)
	batch.PUT("/update", middlewares.HasPermission("batch:update"), c.UpdateBatch)
	batch.DELETE("/delete", middlewares.HasPermission("batch:delete"), c.DeleteBatch)
}

// GetBatchList 获取批次列表
func (c *BatchController) GetBatchList(ctx *gin.Context) {
	var req batchApi.BatchQueryReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		zap.L().Error("绑定查询参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}
	list, total, err := c.batchService.GetBatchList(ctx, &req)
	if err != nil {
		zap.L().Error("查询批次列表失败", zap.Error(err))
		baizeContext.Waring(ctx, "查询失败")
		return
	}
	baizeContext.SuccessData(ctx, gin.H{"rows": list, "total": total})
}

// GetBatchById 根据ID获取批次
func (c *BatchController) GetBatchById(ctx *gin.Context) {
	batchIdStr := ctx.Param("batchId")
	batchId, err := strconv.ParseInt(batchIdStr, 10, 64)
	if err != nil {
		zap.L().Error("解析批次ID失败", zap.String("batchId", batchIdStr), zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}
	batch, err := c.batchService.GetBatchById(ctx, batchId)
	if err != nil {
		zap.L().Error("查询批次失败", zap.Int64("batchId", batchId), zap.Error(err))
		baizeContext.Waring(ctx, "查询失败")
		return
	}
	if batch == nil {
		baizeContext.Waring(ctx, "批次不存在")
		return
	}
	baizeContext.SuccessData(ctx, batch)
}

// CreateBatch 创建批次
func (c *BatchController) CreateBatch(ctx *gin.Context) {
	var req batchApi.BatchCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.L().Error("绑定创建参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}
	if err := c.batchService.CreateBatch(ctx, &req); err != nil {
		zap.L().Error("创建批次失败", zap.Error(err))
		baizeContext.Waring(ctx, "创建失败")
		return
	}
	baizeContext.Success(ctx)
}

// UpdateBatch 更新批次
func (c *BatchController) UpdateBatch(ctx *gin.Context) {
	var req batchApi.BatchUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.L().Error("绑定更新参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}
	if err := c.batchService.UpdateBatch(ctx, &req); err != nil {
		zap.L().Error("更新批次失败", zap.Error(err))
		baizeContext.Waring(ctx, "更新失败")
		return
	}
	baizeContext.Success(ctx)
}

// DeleteBatch 删除批次
func (c *BatchController) DeleteBatch(ctx *gin.Context) {
	var req batchApi.BatchDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.L().Error("绑定删除参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}
	if err := c.batchService.DeleteBatch(ctx, &req); err != nil {
		zap.L().Error("删除批次失败", zap.Error(err))
		baizeContext.Waring(ctx, "删除失败")
		return
	}
	baizeContext.Success(ctx)
}
