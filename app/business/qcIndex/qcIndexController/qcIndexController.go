package qcIndexController

import (
	"nova-factory-server/app/business/qcIndex/qcIndexApi"
	"nova-factory-server/app/business/qcIndex/qcIndexService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type QcIndexController struct {
	qcIndexService qcIndexService.QcIndexService
}

func NewQcIndexController(qcIndexService qcIndexService.QcIndexService) *QcIndexController {
	return &QcIndexController{
		qcIndexService: qcIndexService,
	}
}

func (c *QcIndexController) PrivateRoutes(router *gin.RouterGroup) {
	qcIndexRouter := router.Group("/qcIndex")
	qcIndexRouter.GET("/list", middlewares.HasPermission("qcIndex:list"), c.List)                     // 获取检测项列表
	qcIndexRouter.POST("/create", middlewares.HasPermission("qcIndex:create"), c.Create)              // 创建检测项
	qcIndexRouter.PUT("/update", middlewares.HasPermission("qcIndex:update"), c.Update)               // 修改检测项
	qcIndexRouter.DELETE("/delete/:index_ids", middlewares.HasPermission("qcIndex:delete"), c.Delete) // 删除检测项
	qcIndexRouter.GET("/get/:index_id", middlewares.HasPermission("qcIndex:get"), c.GetById)          // 根据ID获取检测项
}

// List 获取检测项列表
// @Summary 获取检测项列表
// @Description 获取检测项列表
// @Tags 检测项管理
// @Accept json
// @Produce json
// @Param indexCode query string false "检测项编码"
// @Param indexName query string false "检测项名称"
// @Param indexType query string false "检测项类型"
// @Param qcResultType query string false "质检值类型"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} qcIndexApi.QcIndexListRes
// @Router /qcIndex/list [get]
func (c *QcIndexController) List(ctx *gin.Context) {
	req := new(qcIndexApi.QcIndexListReq)
	err := ctx.ShouldBind(req)
	if err != nil {
		zap.L().Error("解析错误", zap.Error(err))
		baizeContext.SuccessData(ctx, nil)
		return
	}
	res, err := c.qcIndexService.List(ctx, req)
	if err != nil {
		zap.L().Error("检测项列表错误", zap.Error(err))
		baizeContext.SuccessData(ctx, nil)
		return
	}
	baizeContext.SuccessData(ctx, res)
}

// Create 创建检测项
// @Summary 创建检测项
// @Description 创建检测项
// @Tags 检测项管理
// @Accept json
// @Produce json
// @Param qcIndex body qcIndexApi.QcIndexCreateReq true "检测项信息"
// @Success 200 {object} qcIndexApi.QcIndexCreateRes
// @Router /qcIndex/create [post]
func (c *QcIndexController) Create(ctx *gin.Context) {
	req := new(qcIndexApi.QcIndexCreateReq)
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		zap.L().Error("解析错误", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}
	res, err := c.qcIndexService.Create(ctx, req)
	if err != nil {
		zap.L().Error("创建检测项错误", zap.Error(err))
		baizeContext.Waring(ctx, "创建检测项失败")
		return
	}
	baizeContext.SuccessData(ctx, res)
}

// Update 修改检测项
// @Summary 修改检测项
// @Description 修改检测项
// @Tags 检测项管理
// @Accept json
// @Produce json
// @Param qcIndex body qcIndexApi.QcIndexUpdateReq true "检测项信息"
// @Success 200 {object} qcIndexApi.QcIndexUpdateRes
// @Router /qcIndex/update [put]
func (c *QcIndexController) Update(ctx *gin.Context) {
	req := new(qcIndexApi.QcIndexUpdateReq)
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		zap.L().Error("解析错误", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}
	res, err := c.qcIndexService.Update(ctx, req)
	if err != nil {
		zap.L().Error("修改检测项错误", zap.Error(err))
		baizeContext.Waring(ctx, "修改检测项失败")
		return
	}
	baizeContext.SuccessData(ctx, res)
}

// Delete 删除检测项
// @Summary 删除检测项
// @Description 删除检测项
// @Tags 检测项管理
// @Param index_ids path string true "检测项ID列表，多个用逗号分隔"
// @Produce json
// @Success 200 {object} response.ResponseData
// @Router /qcIndex/delete/{index_ids} [delete]
func (c *QcIndexController) Delete(ctx *gin.Context) {
	indexIds := baizeContext.ParamInt64Array(ctx, "index_ids")
	if len(indexIds) == 0 {
		baizeContext.Waring(ctx, "请选择要删除的检测项")
		return
	}
	err := c.qcIndexService.Delete(ctx, indexIds)
	if err != nil {
		zap.L().Error("删除检测项错误", zap.Error(err))
		baizeContext.Waring(ctx, "删除检测项失败")
		return
	}
	baizeContext.Success(ctx)
}

// GetById 根据ID获取检测项
// @Summary 根据ID获取检测项
// @Description 根据ID获取检测项详情
// @Tags 检测项管理
// @Param index_id path int true "检测项ID"
// @Produce json
// @Success 200 {object} qcIndexApi.QcIndexUpdateRes
// @Router /qcIndex/get/{index_id} [get]
func (c *QcIndexController) GetById(ctx *gin.Context) {
	indexId := baizeContext.ParamInt64(ctx, "index_id")
	if indexId == 0 {
		baizeContext.Waring(ctx, "检测项ID不能为空")
		return
	}
	res, err := c.qcIndexService.GetById(ctx, indexId)
	if err != nil {
		zap.L().Error("获取检测项详情错误", zap.Error(err))
		baizeContext.Waring(ctx, "获取检测项详情失败")
		return
	}
	baizeContext.SuccessData(ctx, res)
}
