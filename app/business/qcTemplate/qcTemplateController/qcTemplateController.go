package qcTemplateController

import (
	"strconv"

	"nova-factory-server/app/business/qcTemplate/qcTemplateApi"
	"nova-factory-server/app/business/qcTemplate/qcTemplateService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// QcTemplateController 检测模板控制器
type QcTemplateController struct {
	qcTemplateService qcTemplateService.IQcTemplateService
}

// NewQcTemplateController 创建检测模板控制器
func NewQcTemplateController(qcTemplateService qcTemplateService.IQcTemplateService) *QcTemplateController {
	return &QcTemplateController{
		qcTemplateService: qcTemplateService,
	}
}

// PrivateRoutes 私有路由
func (c *QcTemplateController) PrivateRoutes(router *gin.RouterGroup) {
	qcTemplate := router.Group("/qc/template")
	qcTemplate.GET("/list", middlewares.HasPermission("qc:template:list"), c.GetQcTemplateList)            // 获取检测模板列表
	qcTemplate.GET("/get/:templateId", middlewares.HasPermission("qc:template:get"), c.GetQcTemplateById)  // 根据ID获取检测模板
	qcTemplate.POST("/create", middlewares.HasPermission("qc:template:create"), c.CreateQcTemplate)        // 创建检测模板
	qcTemplate.PUT("/update", middlewares.HasPermission("qc:template:update"), c.UpdateQcTemplate)         // 更新检测模板
	qcTemplate.DELETE("/delete", middlewares.HasPermission("qc:template:delete"), c.DeleteQcTemplateByIds) // 删除检测模板
}

// GetQcTemplateList 获取检测模板列表
// @Summary 获取检测模板列表
// @Description 分页查询检测模板列表
// @Tags 检测模板
// @Accept json
// @Produce json
// @Param pageNum query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param templateCode query string false "检测模板编号"
// @Param templateName query string false "检测模板名称"
// @Param qcTypes query string false "检测种类"
// @Param enableFlag query string false "是否启用"
// @Param createBy query string false "创建者"
// @Router /qc/template/list [get]
func (c *QcTemplateController) GetQcTemplateList(ctx *gin.Context) {
	var req qcTemplateApi.QcTemplateQueryReq
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

	list, total, err := c.qcTemplateService.GetQcTemplateList(ctx, &req)
	if err != nil {
		zap.L().Error("查询检测模板列表失败", zap.Error(err))
		baizeContext.Waring(ctx, "查询失败")
		return
	}

	// 返回结果
	baizeContext.SuccessData(ctx, gin.H{
		"rows":  list,
		"total": total,
	})
}

// GetQcTemplateById 根据ID获取检测模板
// @Summary 根据ID获取检测模板
// @Description 根据ID获取检测模板详情
// @Tags 检测模板
// @Accept json
// @Produce json
// @Param templateId path int true "检测模板ID"
// @Router /qc/template/get/{templateId} [get]
func (c *QcTemplateController) GetQcTemplateById(ctx *gin.Context) {
	templateIdStr := ctx.Param("templateId")
	templateId, err := strconv.ParseInt(templateIdStr, 10, 64)
	if err != nil {
		baizeContext.ParameterError(ctx)
		return
	}

	template, err := c.qcTemplateService.GetQcTemplateById(ctx, templateId)
	if err != nil {
		zap.L().Error("查询检测模板失败", zap.Error(err))
		baizeContext.Waring(ctx, "查询失败")
		return
	}

	if template == nil {
		baizeContext.Waring(ctx, "检测模板不存在")
		return
	}

	// 返回结果
	baizeContext.SuccessData(ctx, template)
}

// CreateQcTemplate 创建检测模板
// @Summary 创建检测模板
// @Description 创建新的检测模板
// @Tags 检测模板
// @Accept json
// @Produce json
// @Param template body qcTemplateApi.QcTemplateCreateReq true "检测模板信息"
// @Router /qc/template/create [post]
func (c *QcTemplateController) CreateQcTemplate(ctx *gin.Context) {
	var req qcTemplateApi.QcTemplateCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.L().Error("绑定创建参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	err := c.qcTemplateService.CreateQcTemplate(ctx, &req)
	if err != nil {
		zap.L().Error("创建检测模板失败", zap.Error(err))
		baizeContext.Waring(ctx, "创建失败")
		return
	}

	// 返回结果
	baizeContext.Success(ctx)
}

// UpdateQcTemplate 更新检测模板
// @Summary 更新检测模板
// @Description 更新检测模板信息
// @Tags 检测模板
// @Accept json
// @Produce json
// @Param template body qcTemplateApi.QcTemplateUpdateReq true "检测模板信息"
// @Router /qc/template/update [put]
func (c *QcTemplateController) UpdateQcTemplate(ctx *gin.Context) {
	var req qcTemplateApi.QcTemplateUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.L().Error("绑定更新参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	err := c.qcTemplateService.UpdateQcTemplate(ctx, &req)
	if err != nil {
		zap.L().Error("更新检测模板失败", zap.Error(err))
		baizeContext.Waring(ctx, "更新失败")
		return
	}

	// 返回结果
	baizeContext.Success(ctx)
}

// DeleteQcTemplateByIds 批量删除检测模板
// @Summary 批量删除检测模板
// @Description 批量删除检测模板
// @Tags 检测模板
// @Accept json
// @Produce json
// @Param templateIds body []int64 true "检测模板ID列表"
// @Router /qc/template/delete [delete]
func (c *QcTemplateController) DeleteQcTemplateByIds(ctx *gin.Context) {
	var templateIds []int64
	if err := ctx.ShouldBindJSON(&templateIds); err != nil {
		zap.L().Error("绑定删除参数失败", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	err := c.qcTemplateService.DeleteQcTemplateByIds(ctx, templateIds)
	if err != nil {
		zap.L().Error("删除检测模板失败", zap.Error(err))
		baizeContext.Waring(ctx, "删除失败")
		return
	}

	// 返回结果
	baizeContext.Success(ctx)
}
