package mastercontroller

import (
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Product ERP 产品控制器
type Product struct {
	service masterservice.IProductService
}

// NewProduct 创建 ERP 产品控制器。
func NewProduct(service masterservice.IProductService) *Product {
	return &Product{service: service}
}

// PrivateRoutes 注册 ERP 产品私有路由。
func (o *Product) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/master/product")
	group.GET("/list", middlewares.HasPermission("erp:master:product:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:master:product:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:master:product:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:master:product:remove"), o.Delete)
	group.POST("/vector/generate/:id", middlewares.HasPermission("erp:master:product:vector:generate"), o.Generate)
	group.POST("/vector/generate/all", middlewares.HasPermission("erp:master:product:vector:generate:all"), o.GenerateAll)
	group.GET("/vector/generate/all/progress/:taskId",
		middlewares.HasPermission("erp:master:product:vector:generate:all:progress"), o.GetGenerateAllProgress)
	group.POST("/vector/search", middlewares.HasPermission("erp:master:product:vector:search"), o.Search)
	group.POST("/vector/search/batch", middlewares.HasPermission("erp:master:product:vector:batch_search"), o.BatchSearch)
}

// List 查询 ERP 产品列表。
// @Summary 查询 ERP 产品列表
// @Description 按条件分页查询 ERP 产品列表
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param object query mastermodels.ProductQuery true "ERP 产品查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/master/product/list [get]
func (o *Product) List(c *gin.Context) {
	req := new(mastermodels.ProductQuery)
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

// GetByID 查询 ERP 产品详情。
// @Summary 查询 ERP 产品详情
// @Description 根据ID查询 ERP 产品详情
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/master/product/query/{id} [get]
func (o *Product) GetByID(c *gin.Context) {
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

// Set 新增或修改 ERP 产品。
// @Summary 新增或修改 ERP 产品
// @Description 新增或修改 ERP 产品
// @Tags ERP/基础资料
// @Security BearerAuth
// @Accept application/json
// @Param body body mastermodels.ProductUpsert true "ERP 产品参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/master/product/set [post]
func (o *Product) Set(c *gin.Context) {
	req := new(mastermodels.ProductUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *mastermodels.Product
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

// Delete 删除 ERP 产品。
// @Summary 删除 ERP 产品
// @Description 根据ID删除 ERP 产品，多个ID用逗号分隔
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/master/product/remove/{ids} [delete]
func (o *Product) Delete(c *gin.Context) {
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

// Generate 生成产品向量。
// @Summary 生成产品向量
// @Description 根据产品ID生成产品向量
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param id path int true "产品ID"
// @Param body body mastermodels.ProductGenVectorReq true "产品向量生成参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "生成成功"
// @Router /erp/master/product/vector/generate/{id} [post]
func (o *Product) Generate(c *gin.Context) {
	req := new(mastermodels.ProductGenVectorReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.ID <= 0 {
		req.ID = baizeContext.ParamInt64(c, "id")
	}
	if req.ID <= 0 || req.Embedding == nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.GenerateVector(c, req)
	if err != nil {
		zap.L().Error("generate product vector fail", zap.Int64("id", req.ID), zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GenerateAll 全量生成启用产品向量。
// @Summary 全量生成启用产品向量
// @Description 异步将全部启用产品写入向量数据库，并将任务进度记录到 Redis
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param body body mastermodels.ProductGenAllVectorReq true "全量生成产品向量参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "任务已启动"
// @Router /erp/master/product/vector/generate/all [post]
func (o *Product) GenerateAll(c *gin.Context) {
	req := new(mastermodels.ProductGenAllVectorReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.Embedding == nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.GenerateAllVectors(c, req)
	if err != nil {
		zap.L().Error("generate all product vectors fail", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetGenerateAllProgress 查询全量生成产品向量任务进度。
// @Summary 查询全量生成产品向量任务进度
// @Description 根据任务ID读取 Redis 中的全量生成任务进度
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param taskId path string true "任务ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/master/product/vector/generate/all/progress/{taskId} [get]
func (o *Product) GetGenerateAllProgress(c *gin.Context) {
	taskID := c.Param("taskId")
	if taskID == "" {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.GetGenerateAllVectorsProgress(c, taskID)
	if err != nil {
		zap.L().Error("get product vector progress fail", zap.String("taskId", taskID), zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Search 近似搜索产品向量。
// @Summary 近似搜索产品向量
// @Description 根据检索文本生成查询向量并近似搜索产品数据
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param body body mastermodels.ProductVectorSearchReq true "产品向量搜索参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "搜索成功"
// @Router /erp/master/product/vector/search [post]
func (o *Product) Search(c *gin.Context) {
	req := new(mastermodels.ProductVectorSearchReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.Embedding == nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.SearchVector(c, req)
	if err != nil {
		zap.L().Error("search product vector fail", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// BatchSearch 批量近似搜索产品向量。
// @Summary 批量近似搜索产品向量
// @Description 根据多条检索文本批量生成查询向量并近似搜索产品数据
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param body body mastermodels.ProductVectorBatchSearchReq true "产品批量向量搜索参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "搜索成功"
// @Router /erp/master/product/vector/search/batch [post]
func (o *Product) BatchSearch(c *gin.Context) {
	req := new(mastermodels.ProductVectorBatchSearchReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.Embedding == nil || len(req.Queries) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.BatchSearchVector(c, req)
	if err != nil {
		zap.L().Error("batch search product vector fail", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}
