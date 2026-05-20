package aidatasetcontroller

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	server "nova-factory-server/app/utils/gin_mcp"
)

type DataSetAuth struct {
	// 权限
	permissionService aidatasetservice.IDatasetRolePermissionService
}

func NewDataSetAuth(permissionService aidatasetservice.IDatasetRolePermissionService) *DataSetAuth {
	return &DataSetAuth{
		permissionService: permissionService,
	}
}

func (d *DataSetAuth) PrivateRoutes(router *gin.RouterGroup) {
	ai := router.Group("/ai/dataset")
	ai.POST("/auth", middlewares.HasPermission("ai:dataset:auth"), d.Auth) // 创建知识库
}

// PrivateMcpRoutes 注册mcp服务
func (d *DataSetAuth) PrivateMcpRoutes(server *server.GinMCP) {
	server.RegisterSchema("POST", "/ai/dataset/auth", nil, nil) // 创建知识库
}

// Auth 获取当前用户知识库授权信息
// @Summary 获取知识库授权信息.
// @Description 在调用 MCP 工具 `ragflow_retrieval` 前先使用此工具。它会根据用户访问权限返回访问所需的 `dataset_ids` 和 `document_ids`。
// @Tags 工业智能体/知识库管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseData "成功"
// @Router /ai/dataset/auth [post]
func (d *DataSetAuth) Auth(c *gin.Context) {
	data, err := d.permissionService.GetDatasetData(c)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}

	if len(data.DatasetUuIDs) == 0 {
		baizeContext.Waring(c, "用户无知识库权限")
		return
	}

	var output aidatasetmodels.RagflowAuthOutput
	output.DocumentIDs = data.DocumentUuIDs
	output.DatasetIDs = data.DatasetUuIDs
	output.Configured = true
	baizeContext.SuccessData(c, output)
	return
}
