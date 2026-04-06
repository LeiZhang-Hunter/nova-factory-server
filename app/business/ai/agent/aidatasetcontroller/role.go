package aidatasetcontroller

import (
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Role struct {
	service aidatasetservice.IDatasetRolePermissionService
}

// NewRole 知识库权限控制器构造函数。
func NewRole(service aidatasetservice.IDatasetRolePermissionService) *Role {
	return &Role{
		service: service,
	}
}

// PrivateRoutes 注册知识库/文档权限相关私有路由。
func (r *Role) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/dataset/role/permission")
	group.GET("/list", middlewares.HasPermission("ai:dataset:role:permission:list"), r.List)
	group.POST("/set", middlewares.HasPermission("ai:dataset:role:permission:set"), r.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("ai:dataset:role:permission:remove"), r.Remove)
}

// List 查询知识库/文档角色权限列表
// @Summary 查询知识库/文档角色权限列表
// @Description 查询知识库/文档角色权限列表
// @Tags 工业智能体/知识库权限
// @Param object query aidatasetmodels.DatasetRolePermissionQuery true "权限列表查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /ai/dataset/role/permission/list [get]
func (r *Role) List(c *gin.Context) {
	req := new(aidatasetmodels.DatasetRolePermissionQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := r.service.List(c, req)
	if err != nil {
		zap.L().Error("list role permission error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set 新增或修改知识库/文档角色权限
// @Summary 新增或修改知识库/文档角色权限
// @Description 传入id时修改，不传id时新增。documentId 为0表示对整个知识库授权
// @Tags 工业智能体/知识库权限
// @Param object body aidatasetmodels.SetDatasetRolePermission true "权限设置参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "操作成功"
// @Router /ai/dataset/role/permission/set [post]
func (r *Role) Set(c *gin.Context) {
	req := new(aidatasetmodels.SetDatasetRolePermission)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := r.service.Set(c, req)
	if err != nil {
		zap.L().Error("set role permission error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Remove 删除知识库/文档角色权限
// @Summary 删除知识库/文档角色权限
// @Description 删除知识库/文档角色权限（软删除）
// @Tags 工业智能体/知识库权限
// @Param ids path string true "权限ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /ai/dataset/role/permission/remove/{ids} [delete]
func (r *Role) Remove(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := r.service.Remove(c, ids); err != nil {
		zap.L().Error("remove role permission error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
