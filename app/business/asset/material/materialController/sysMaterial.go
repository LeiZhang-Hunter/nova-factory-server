package materialController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/material/materialModels"
	"nova-factory-server/app/business/asset/material/materialService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

// MaterialInfo 物料管理
type MaterialInfo struct {
	iMaterialService materialService.IMaterialService
}

func NewMaterialInfo(iMaterialService materialService.IMaterialService) *MaterialInfo {
	return &MaterialInfo{
		iMaterialService: iMaterialService,
	}
}

func (mi *MaterialInfo) PrivateRoutes(router *gin.RouterGroup) {
	material := router.Group("/asset/material")
	material.GET("/list", middlewares.HasPermission("asset:material"), mi.GetMaterialInfoList)                  // 物料列表
	material.POST("/set", middlewares.HasPermission("asset:material:set"), mi.SetMaterialInfo)                  // 设置物料信息
	material.DELETE("/:materialIds", middlewares.HasPermission("asset:material:remove"), mi.RemoveMaterialInfo) //删除物料

	material.GET("/inbound/list", middlewares.HasPermission("asset:material:inbound"), mi.GetMaterialInBoundList) // 入库列表
	material.POST("/inbound/add", middlewares.HasPermission("asset:material:inbound:add"), mi.AddMaterialInBound) // 入库登记

	material.GET("/outbound/list", middlewares.HasPermission("asset:material:outbound"), mi.GetMaterialOutBoundList) // 出库登记
	material.POST("/outbound/add", middlewares.HasPermission("asset:material:outbound:add"), mi.AddMaterialOutBound) // 出库登记
}

// GetMaterialInfoList 物料管理列表
// @Summary 物料管理列表
// @Description 物料管理列表
// @Tags 物料管理
// @Param  object query materialModels.MaterialListReq true "物料管理列表请求参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/material/list [get]
func (mi *MaterialInfo) GetMaterialInfoList(c *gin.Context) {
	req := new(materialModels.MaterialListReq)
	err := c.ShouldBindQuery(req)
	list, err := mi.iMaterialService.SelectMaterialList(c, req)
	if err != nil {
		zap.L().Error("读取设备分组失败", zap.Error(err))
	}
	baizeContext.SuccessData(c, list)
}

// SetMaterialInfo 登记物料
// @Summary 登记物料
// @Description 登记物料
// @Tags 物料管理
// @Param  object body materialModels.MaterialInfo true "物料参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置成功"
// @Router /asset/material/set [post]
func (mi *MaterialInfo) SetMaterialInfo(c *gin.Context) {
	info := new(materialModels.MaterialInfo)
	err := c.ShouldBindJSON(info)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if info.MaterialId > 0 {
		vo, err := mi.iMaterialService.UpdateMaterial(c, info)
		if err != nil {
			zap.L().Error("更新设备数据失败", zap.Error(err))
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, vo)
	} else {
		vo, err := mi.iMaterialService.InsertMaterial(c, info)
		if err != nil {
			zap.L().Error("插入设备数据失败", zap.Error(err))
			baizeContext.Waring(c, "插入设备数据失败")
			return
		}
		baizeContext.SuccessData(c, vo)
	}
}

// RemoveMaterialInfo 删除物料
// @Summary 删除物料
// @Description 删除物料
// @Tags 物料管理
// @Param  materialIds path string true "materialIds"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /asset/material/{materialIds}  [delete]
func (mi *MaterialInfo) RemoveMaterialInfo(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "materialIds")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		return
	}
	err := mi.iMaterialService.DeleteByMaterialIds(c, ids)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}

	baizeContext.Success(c)
}

// GetMaterialInBoundList 物料入库列表
// @Summary 物料入库列表
// @Description 物料入库列表
// @Tags 物料管理
// @Param  object query materialModels.InboundListReq true "物料入库管理查询条件"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /asset/material/inbound/list  [get]
func (mi *MaterialInfo) GetMaterialInBoundList(c *gin.Context) {
	req := new(materialModels.InboundListReq)
	err := c.ShouldBindQuery(req)
	list, err := mi.iMaterialService.InboundList(c, req)
	if err != nil {
		zap.L().Error("读取设备分组失败", zap.Error(err))
		baizeContext.SuccessData(c, &materialModels.InboundListValue{
			Rows:  make([]*materialModels.InboundValue, 0),
			Total: 0,
		})
		return
	}
	baizeContext.SuccessData(c, list)
}

// AddMaterialInBound 登记物料入库
// @Summary 登记物料入库
// @Description 登记物料入库
// @Tags 物料管理
// @Param  object body materialModels.InboundInfo true "登记物料入库参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置成功"
// @Router /asset/material/inbound/add [post]
func (mi *MaterialInfo) AddMaterialInBound(c *gin.Context) {
	info := new(materialModels.InboundInfo)
	err := c.ShouldBindJSON(info)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	if info.Number <= 0 {
		baizeContext.Waring(c, "入库数量不能为0")
		return
	}

	vo, err := mi.iMaterialService.GetByMaterialId(c, int64(info.MaterialId))
	if err != nil {
		zap.L().Error("查询物料失败", zap.Error(err))
		baizeContext.Waring(c, "查询物料失败")
		return
	}

	if vo == nil {
		baizeContext.Waring(c, "物料不存在")
		return
	}

	inbound, err := mi.iMaterialService.Inbound(c, info)
	if err != nil {
		zap.L().Error("入库失败", zap.Error(err))
		baizeContext.Waring(c, "入库失败")
		return
	}

	baizeContext.SuccessData(c, inbound)

}

// GetMaterialOutBoundList 物料出库列表
// @Summary 物料出库列表
// @Description 物料出库列表
// @Tags 物料管理
// @Param  object query materialModels.OutboundListReq true "物料出库列表查询条件"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /asset/material/outbound/list  [get]
func (mi *MaterialInfo) GetMaterialOutBoundList(c *gin.Context) {
	req := new(materialModels.OutboundListReq)
	err := c.ShouldBindQuery(req)
	list, err := mi.iMaterialService.OutboundList(c, req)
	if err != nil {
		zap.L().Error("读取设备分组失败", zap.Error(err))
		baizeContext.SuccessData(c, &materialModels.OutboundListValue{
			Rows:  make([]*materialModels.OutboundValue, 0),
			Total: 0,
		})
		return
	}
	baizeContext.SuccessData(c, list)
}

// AddMaterialOutBound 登记物料出库
// @Summary 登记物料出库
// @Description 登记物料出库
// @Tags 物料管理
// @Param  object body materialModels.OutboundInfo true "登记物料出库参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置成功"
// @Router /asset/material/outbound/add [post]
func (mi *MaterialInfo) AddMaterialOutBound(c *gin.Context) {
	info := new(materialModels.OutboundInfo)
	err := c.ShouldBindJSON(info)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	if info.Number <= 0 {
		baizeContext.Waring(c, "入库数量不能为0")
		return
	}

	vo, err := mi.iMaterialService.GetByMaterialId(c, int64(info.MaterialId))
	if err != nil {
		zap.L().Error("查询物料失败", zap.Error(err))
		baizeContext.Waring(c, "查询物料失败")
		return
	}

	if vo == nil {
		baizeContext.Waring(c, "物料不存在")
		return
	}

	if info.Reason == "" {
		info.Reason = "无出库原因"
	}

	outbound, err := mi.iMaterialService.Outbound(c, info)
	if err != nil {
		zap.L().Error("出库失败", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}

	baizeContext.SuccessData(c, outbound)
}
