package defectController

import (
	"nova-factory-server/app/business/defect/defectApi"
	"nova-factory-server/app/business/defect/defectService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DefectController struct {
	defectService defectService.DefectService
}

func NewDefectController(defectService defectService.DefectService) *DefectController {
	return &DefectController{
		defectService: defectService,
	}
}

func (c *DefectController) PrivateRoutes(router *gin.RouterGroup) {
	defectRouter := router.Group("/defect")
	defectRouter.GET("/list", middlewares.HasPermission("defect:list"), c.List)                      // 获取缺陷列表
	defectRouter.POST("/create", middlewares.HasPermission("defect:create"), c.Create)               // 创建缺陷
	defectRouter.PUT("/update", middlewares.HasPermission("defect:update"), c.Update)                // 修改缺陷
	defectRouter.DELETE("/delete/:defect_ids", middlewares.HasPermission("defect:delete"), c.Delete) // 删除缺陷
	defectRouter.GET("/get/:defect_id", middlewares.HasPermission("defect:get"), c.GetById)          // 根据ID获取缺陷
}

// List 获取缺陷列表
// @Summary 获取缺陷列表
// @Description 获取缺陷列表
// @Tags 缺陷管理
// @Accept json
// @Produce json
// @Param defectCode query string false "缺陷编码"
// @Param defectName query string false "缺陷名称"
// @Param indexType query string false "指标类型"
// @Param defectLevel query string false "缺陷等级"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Router /defect/list [get]
func (c *DefectController) List(ctx *gin.Context) {
	req := new(defectApi.DefectListReq)
	err := ctx.ShouldBind(req)
	if err != nil {
		zap.L().Error("解析错误", zap.Error(err))
		baizeContext.SuccessData(ctx, nil)
		return
	}
	res, err := c.defectService.List(ctx, req)
	if err != nil {
		zap.L().Error("缺陷列表错误", zap.Error(err))
		baizeContext.SuccessData(ctx, nil)
		return
	}
	baizeContext.SuccessData(ctx, res)
}

// Create 创建缺陷
// @Summary 创建缺陷
// @Description 创建缺陷
// @Tags 缺陷管理
// @Accept json
// @Produce json
// @Param defect body defectApi.DefectCreateReq true "缺陷信息"
// @Success 200 {object} defectApi.DefectCreateRes
// @Router /defect/create [post]
func (c *DefectController) Create(ctx *gin.Context) {
	req := new(defectApi.DefectCreateReq)
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		zap.L().Error("解析错误", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}
	res, err := c.defectService.Create(ctx, req)
	if err != nil {
		zap.L().Error("创建缺陷错误", zap.Error(err))
		baizeContext.Waring(ctx, "创建缺陷失败")
		return
	}
	baizeContext.SuccessData(ctx, res)
}

// Update 修改缺陷
// @Summary 修改缺陷
// @Description 修改缺陷
// @Tags 缺陷管理
// @Accept json
// @Produce json
// @Param defect body defectApi.DefectUpdateReq true "缺陷信息"
// @Success 200 {object} defectApi.DefectUpdateRes
// @Router /defect/update [put]
func (c *DefectController) Update(ctx *gin.Context) {
	req := new(defectApi.DefectUpdateReq)
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		zap.L().Error("解析错误", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}
	res, err := c.defectService.Update(ctx, req)
	if err != nil {
		zap.L().Error("修改缺陷错误", zap.Error(err))
		baizeContext.Waring(ctx, "修改缺陷失败")
		return
	}
	baizeContext.SuccessData(ctx, res)
}

// Delete 删除缺陷
// @Summary 删除缺陷
// @Description 删除缺陷
// @Tags 缺陷管理
// @Param defect_ids path string true "缺陷ID列表，多个用逗号分隔"
// @Produce json
// @Success 200 {object} response.ResponseData
// @Router /defect/delete/{defect_ids} [delete]
func (c *DefectController) Delete(ctx *gin.Context) {
	defectIds := baizeContext.ParamInt64Array(ctx, "defect_ids")
	if len(defectIds) == 0 {
		baizeContext.Waring(ctx, "请选择要删除的缺陷")
		return
	}
	err := c.defectService.Delete(ctx, defectIds)
	if err != nil {
		zap.L().Error("删除缺陷错误", zap.Error(err))
		baizeContext.Waring(ctx, "删除缺陷失败")
		return
	}
	baizeContext.Success(ctx)
}

// GetById 根据ID获取缺陷
// @Summary 根据ID获取缺陷
// @Description 根据ID获取缺陷详情
// @Tags 缺陷管理
// @Param defect_id path int true "缺陷ID"
// @Produce json
// @Success 200 {object} defectApi.DefectUpdateRes
// @Router /defect/get/{defect_id} [get]
func (c *DefectController) GetById(ctx *gin.Context) {
	defectId := baizeContext.ParamInt64(ctx, "defect_id")
	if defectId == 0 {
		baizeContext.Waring(ctx, "缺陷ID不能为空")
		return
	}
	res, err := c.defectService.GetById(ctx, defectId)
	if err != nil {
		zap.L().Error("获取缺陷详情错误", zap.Error(err))
		baizeContext.Waring(ctx, "获取缺陷详情失败")
		return
	}
	baizeContext.SuccessData(ctx, res)
}
