package resourceController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/resource/resourceModels"
	"nova-factory-server/app/business/asset/resource/resourceService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type ResourceFile struct {
	service resourceService.IResourceFileService
}

func NewResourceFile(service resourceService.IResourceFileService) *ResourceFile {
	return &ResourceFile{service: service}
}

func (rf *ResourceFile) PrivateRoutes(router *gin.RouterGroup) {
	resource := router.Group("/asset/resource")
	resource.GET("/list", middlewares.HasPermission("asset:resource"), rf.List)                    // 资料列表
	resource.POST("/set", middlewares.HasPermission("asset:resource:set"), rf.Set)                 // 设置资料信息
	resource.DELETE("/remove/:ids", middlewares.HasPermission("asset:resource:remove"), rf.Remove) //删除资料
}

// Set 登记资料
// @Summary 登记资料
// @Description 登记资料
// @Tags 资料管理
// @Param  object body resourceModels.SysResourceFileDML true "资料参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置成功"
// @Router /asset/resource/set [post]
func (rf *ResourceFile) Set(c *gin.Context) {
	req := new(resourceModels.SysResourceFileDML)
	err := c.ShouldBindJSON(req)
	if err != nil {
		zap.L().Error("param error", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	if req.ResourceId <= 0 {
		resource, err := rf.service.InsertResource(c, req)
		if err != nil {
			baizeContext.Waring(c, "创建资源失败")
			return
		}
		baizeContext.SuccessData(c, resource)
	} else {
		resource, err := rf.service.UpdateResource(c, req)
		if err != nil {
			baizeContext.Waring(c, "更新资源失败")
			return
		}
		baizeContext.SuccessData(c, resource)
	}
}

// List 资料管理列表
// @Summary 资料管理列表
// @Description 资料管理列表
// @Tags 资料管理
// @Param  object query resourceModels.SysResourceFileDQL true "资料管理列表请求参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/resource/list [get]
func (rf *ResourceFile) List(c *gin.Context) {
	query := new(resourceModels.SysResourceFileDQL)
	err := c.ShouldBindQuery(query)
	if err != nil {
		baizeContext.ParameterError(c)
		zap.L().Error("param error", zap.Error(err))
		return
	}
	list, err := rf.service.List(c, query)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// Remove 删除资料
// @Summary 删除资料
// @Description 删除资料
// @Tags 资料管理
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /asset/resource/remove/{ids}  [delete]
func (rf *ResourceFile) Remove(c *gin.Context) {
	ids := baizeContext.ParamStringArray(c, "ids")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		return
	}
	err := rf.service.Remove(c, ids)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}

	baizeContext.Success(c)
}
