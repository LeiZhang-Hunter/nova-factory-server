package qcIpqcController

import (
	"nova-factory-server/app/business/qcIpqc/qcIpqcApi"
	"nova-factory-server/app/business/qcIpqc/qcIpqcService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type QcIpqcController struct {
	service qcIpqcService.IQcIpqcService
}

func NewQcIpqcController(service qcIpqcService.IQcIpqcService) *QcIpqcController {
	return &QcIpqcController{service: service}
}

func (c *QcIpqcController) PrivateRoutes(router *gin.RouterGroup) {
	qcIpqc := router.Group("/qc/ipqc")
	qcIpqc.GET("/list", middlewares.HasPermission("qc:ipqc:list"), c.GetQcIpqcList)
	qcIpqc.GET("/get/:ipqcId", middlewares.HasPermission("qc:ipqc:get"), c.GetQcIpqcById)
	qcIpqc.POST("/create", middlewares.HasPermission("qc:ipqc:create"), c.CreateQcIpqc)
	qcIpqc.PUT("/update", middlewares.HasPermission("qc:ipqc:update"), c.UpdateQcIpqc)
	qcIpqc.DELETE("/delete", middlewares.HasPermission("qc:ipqc:delete"), c.DeleteQcIpqcByIds)
}

// GetQcIpqcList 获取过程检验单列表
// @Summary 获取过程检验单列表
// @Description 分页查询过程检验单列表
// @Tags 过程检验单
// @Accept json
// @Produce json
// @Param pageNum query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param ipqcCode query string false "检验单编号"
// @Param ipqcName query string false "检验单名称"
// @Param ipqcType query string false "检验类型"
// @Param templateId query int64 false "检验模板ID"
// @Param status query string false "单据状态"
// @Router /qc/ipqc/list [get]
func (c *QcIpqcController) GetQcIpqcList(ctx *gin.Context) {
	var req qcIpqcApi.QcIpqcQueryReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		zap.L().Error("绑定查询参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	list, total, err := c.service.GetQcIpqcList(ctx, &req)
	if err != nil {
		zap.L().Error("查询过程检验单列表失败", zap.Error(err))
		baizeContext.Waring(ctx, "查询失败")
		return
	}
	baizeContext.SuccessData(ctx, gin.H{"rows": list, "total": total})
}

// GetQcIpqcById 根据ID获取过程检验单
// @Summary 根据ID获取过程检验单
// @Description 根据ID获取过程检验单详情
// @Tags 过程检验单
// @Accept json
// @Produce json
// @Param ipqcId path int true "过程检验单ID"
// @Router /qc/ipqc/get/{ipqcId} [get]
func (c *QcIpqcController) GetQcIpqcById(ctx *gin.Context) {
	ipqcIdStr := ctx.Param("ipqcId")
	ipqcId, err := strconv.ParseInt(ipqcIdStr, 10, 64)
	if err != nil {
		baizeContext.ParameterError(ctx)
		return
	}
	ipqc, err := c.service.GetQcIpqcById(ctx, ipqcId)
	if err != nil {
		zap.L().Error("查询过程检验单失败", zap.Error(err))
		baizeContext.Waring(ctx, "查询失败")
		return
	}
	if ipqc == nil {
		baizeContext.Waring(ctx, "过程检验单不存在")
		return
	}
	baizeContext.SuccessData(ctx, ipqc)
}

// CreateQcIpqc 创建过程检验单
// @Summary 创建过程检验单
// @Description 创建新的过程检验单
// @Tags 过程检验单
// @Accept json
// @Produce json
// @Param ipqc body qcIpqcApi.QcIpqcCreateReq true "过程检验单信息"
// @Router /qc/ipqc/create [post]
func (c *QcIpqcController) CreateQcIpqc(ctx *gin.Context) {
	var req qcIpqcApi.QcIpqcCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.L().Error("绑定创建参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}
	err := c.service.CreateQcIpqc(ctx, &req)
	if err != nil {
		zap.L().Error("创建过程检验单失败", zap.Error(err))
		baizeContext.Waring(ctx, "创建失败")
		return
	}
	baizeContext.Success(ctx)
}

// UpdateQcIpqc 更新过程检验单
// @Summary 更新过程检验单
// @Description 更新过程检验单信息
// @Tags 过程检验单
// @Accept json
// @Produce json
// @Param ipqc body qcIpqcApi.QcIpqcUpdateReq true "过程检验单信息"
// @Router /qc/ipqc/update [put]
func (c *QcIpqcController) UpdateQcIpqc(ctx *gin.Context) {
	var req qcIpqcApi.QcIpqcUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.L().Error("绑定更新参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}
	err := c.service.UpdateQcIpqc(ctx, &req)
	if err != nil {
		zap.L().Error("更新过程检验单失败", zap.Error(err))
		baizeContext.Waring(ctx, "更新失败")
		return
	}
	baizeContext.Success(ctx)
}

// DeleteQcIpqcByIds 批量删除过程检验单
// @Summary 批量删除过程检验单
// @Description 批量删除过程检验单
// @Tags 过程检验单
// @Accept json
// @Produce json
// @Param ids body []int64 true "过程检验单ID列表"
// @Router /qc/ipqc/delete [delete]
func (c *QcIpqcController) DeleteQcIpqcByIds(ctx *gin.Context) {
	var ids []int64
	if err := ctx.ShouldBindJSON(&ids); err != nil {
		zap.L().Error("绑定删除参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}
	err := c.service.DeleteQcIpqcByIds(ctx, ids)
	if err != nil {
		zap.L().Error("删除过程检验单失败", zap.Error(err))
		baizeContext.Waring(ctx, "删除失败")
		return
	}
	baizeContext.Success(ctx)
}
