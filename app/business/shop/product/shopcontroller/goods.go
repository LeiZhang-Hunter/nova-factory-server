package shopcontroller

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Goods struct {
	service shopservice.IShopGoodsService
}

func NewGoods(service shopservice.IShopGoodsService) *Goods {
	return &Goods{service: service}
}

// PrivateRoutes 注册商品路由
func (s *Goods) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/goods")
	group.GET("/list", middlewares.HasPermission("shop:goods:list"), s.List)
	group.GET("/export/csv", middlewares.HasPermission("shop:goods:export:csv"), s.ExportCSV)
	group.GET("/info/:id", middlewares.HasPermission("shop:goods:info"), s.GetByID)
	group.POST("/add", middlewares.HasPermission("shop:goods:add"), s.Create)
	group.PUT("/edit", middlewares.HasPermission("shop:goods:edit"), s.Update)
	group.DELETE("/remove/:ids", middlewares.HasPermission("shop:goods:remove"), s.Delete)
	group.POST("/vector/generate/:id", middlewares.HasPermission("shop:goods:vector:generate"), s.Generate)
	group.POST("/vector/search", middlewares.HasPermission("shop:goods:vector:search"), s.Search)
}

// PublicRoutes 导入接口注册
func (s *Goods) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/goods")
	group.Any("/import", s.Import)
}

// List 获取商品列表
// @Summary 获取商品列表
// @Description 获取商品列表
// @Tags 商城/商品管理
// @Param object query shopmodels.GoodsQuery true "商品查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/goods/list [get]
func (s *Goods) List(c *gin.Context) {
	req := new(shopmodels.GoodsQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetByID 获取商品详情
// @Summary 获取商品详情
// @Description 根据ID获取商品详情
// @Tags 商城/商品管理
// @Param id path int true "商品ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/goods/{id} [get]
func (s *Goods) GetByID(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.GetByID(c, id)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Create 新增商品
// @Summary 新增商品
// @Description 新增商品
// @Tags 商城/商品管理
// @Param object body shopmodels.GoodsUpsert true "商品新增参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "新增成功"
// @Router /shop/goods [post]
func (s *Goods) Create(c *gin.Context) {
	req := new(shopmodels.GoodsUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.Create(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Update 修改商品
// @Summary 修改商品
// @Description 修改商品
// @Tags 商城/商品管理
// @Param object body shopmodels.GoodsUpsert true "商品修改参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "修改成功"
// @Router /shop/goods [put]
func (s *Goods) Update(c *gin.Context) {
	req := new(shopmodels.GoodsUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.ID == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.Update(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Delete 删除商品
// @Summary 删除商品
// @Description 根据ID删除商品
// @Tags 商城/商品管理
// @Param ids path string true "商品ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/goods/{ids} [delete]
func (s *Goods) Delete(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := s.service.DeleteByIDs(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// Import 导入商品
// @Summary 导入商品
// @Description 导入商品
// @Tags 商城/商品管理
// @Param object body shopmodels.GoodsUpsert true "商品新增参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "新增成功"
// @Router /shop/goods/import [post]
func (s *Goods) Import(c *gin.Context) {
	goods := new(shopmodels.ImportGoodsList)
	err := c.ShouldBindJSON(goods)
	if err != nil {
		zap.L().Error("parse goods fail", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	if err = s.service.Import(c, goods.Records); err != nil {
		zap.L().Error("import goods fail", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// Generate 生成商品向量
// @Summary 生成商品向量
// @Description 根据商品ID生成商品向量
// @Tags 商城/商品管理
// @Param id path int true "商品ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "生成成功"
// @Router /shop/goods/vector/generate/{id} [post]
func (s *Goods) Generate(c *gin.Context) {
	req := new(shopmodels.GenVectorReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.ID <= 0 {
		zap.L().Error("invalid goods id", zap.Int64("id", req.ID))
		baizeContext.ParameterError(c)
		return
	}
	if req.Embedding == nil {
		zap.L().Error("invalid goods embedding")
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.GenerateVector(c, req)
	if err != nil {
		zap.L().Error("generate goods vector fail", zap.Int64("id", req.ID), zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Search 近似搜索商品向量
// @Summary 近似搜索商品向量
// @Description 根据检索文本生成查询向量并近似搜索商品SKU数据
// @Tags 商城/商品管理
// @Param object body shopmodels.GoodsVectorSearchReq true "商品向量搜索参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "搜索成功"
// @Router /shop/goods/vector/search [post]
func (s *Goods) Search(c *gin.Context) {
	req := new(shopmodels.GoodsVectorSearchReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.Embedding == nil {
		zap.L().Error("invalid goods search embedding")
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.SearchVector(c, req)
	if err != nil {
		zap.L().Error("search goods vector fail", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// ExportCSV 导出商品CSV
// @Summary 导出商品CSV
// @Description 按查询条件导出商品CSV文件
// @Tags 商城/商品管理
// @Param object query shopmodels.GoodsQuery true "商品查询参数"
// @Security BearerAuth
// @Produce text/csv
// @Success 200 {file} file "导出成功"
// @Router /shop/goods/export/csv [get]
func (s *Goods) ExportCSV(c *gin.Context) {
	req := new(shopmodels.GoodsQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}

	fileName := fmt.Sprintf("goods_%s.csv", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	c.Header("Cache-Control", "no-store")
	c.Status(http.StatusOK)

	var err error
	if _, err = c.Writer.Write([]byte("\xEF\xBB\xBF")); err != nil {
		zap.L().Error("write csv bom fail", zap.Error(err))
		return
	}

	csvWriter := csv.NewWriter(c.Writer)
	flusher, _ := c.Writer.(http.Flusher)
	if err = s.service.ExportCSV(c, req, csvWriter, func() {
		if flusher != nil {
			flusher.Flush()
		}
	}); err != nil {
		zap.L().Error("export goods csv fail", zap.Error(err))
	}
}
