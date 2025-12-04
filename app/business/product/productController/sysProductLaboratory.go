package productController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/product/productModels"
	"nova-factory-server/app/business/product/productService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Laboratory struct {
	ls productService.ISysProductLaboratoryService
}

func NewLaboratory(ls productService.ISysProductLaboratoryService) *Laboratory {
	return &Laboratory{ls: ls}
}

func (lc *Laboratory) PrivateRoutes(router *gin.RouterGroup) {
	laboratory := router.Group("/product/laboratory")
	laboratory.GET("/list", middlewares.HasPermission("product:laboratory:list"), lc.LaboratoryList)
	laboratory.GET("/info/:id", middlewares.HasPermission("product:laboratory:info"), lc.LaboratoryGetInfo)
	laboratory.POST("/set", middlewares.HasPermission("product:laboratory:set"), lc.Set)
	laboratory.DELETE("/remove/:ids", middlewares.HasPermission("product:laboratory:remove"), lc.Remove)
}

func (lc *Laboratory) PublicRoutes(router *gin.RouterGroup) {
	laboratory := router.Group("/api/v1/product/laboratory")
	laboratory.POST("/set", lc.Set) //登录
}

// LaboratoryList 化验单列表
// @Summary 化验单列表
// @Description 化验单列表
// @Tags 化验单管理
// @Param  object query systemModels.SysProductLaboratoryDQL false "查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData{data=response.ListData{rows=[]systemModels.SysProductLaboratoryVo}} "成功"
// @Router /product/laboratory/list [get]
func (lc *Laboratory) LaboratoryList(c *gin.Context) {
	dql := new(productModels.SysProductLaboratoryDQL)
	_ = c.ShouldBindQuery(dql)
	list, err := lc.ls.SelectLaboratoryList(c, dql)
	if err != nil {
		return
	}
	baizeContext.SuccessData(c, list)
}

// LaboratoryGetInfo 化验单详情
// @Summary 化验单详情
// @Description 化验单详情
// @Tags 化验单管理
// @Param id path string true "主键ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData{data=systemModels.SysProductLaboratoryVo} "成功"
// @Router /product/laboratory/info/{id} [get]
func (lc *Laboratory) LaboratoryGetInfo(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	info, err := lc.ls.SelectLaboratoryById(c, id)
	if err != nil {
		zap.L().Error("get laboratory info", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, info)
}

// Set 保存化验单
// @Summary 保存化验单
// @Description 保存化验单
// @Tags 保存化验单
// @Param  object body systemModels.SysProductLaboratoryVo true "请求体"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /product/laboratory/set [post]
func (lc *Laboratory) Set(c *gin.Context) {
	body := new(productModels.SysProductLaboratoryVo)
	if err := c.ShouldBindJSON(body); err != nil {
		zap.L().Error("设置化验单参数错误", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	info, err := lc.ls.Set(c, body)
	if err != nil {
		zap.L().Error("set laboratory error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, info)
}

// Remove 删除化验单
// @Summary 删除化验单
// @Description 删除化验单
// @Tags 化验单管理
// @Param ids path string true "主键ID,多个用逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /product/laboratory/remove/{ids} [delete]
func (lc *Laboratory) Remove(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	err := lc.ls.DeleteLaboratoryByIds(c, ids)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}
	baizeContext.Success(c)
}
