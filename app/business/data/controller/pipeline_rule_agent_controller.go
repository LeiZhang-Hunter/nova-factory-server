package controller

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/data/agent"
	"nova-factory-server/app/business/data/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"
)

type PipelineRuleAgentController struct {
	agent *agent.Service
	rules service.IPipelineRuleService
}

func NewPipelineRuleAgentController(a *agent.Service, rules service.IPipelineRuleService) *PipelineRuleAgentController {
	return &PipelineRuleAgentController{agent: a, rules: rules}
}

func (c *PipelineRuleAgentController) PrivateRoutes(router *gin.RouterGroup) {
	router.POST("/data/pipeline-rules/generate-from-excel", middlewares.HasPermission("data:service:rule:add"), c.Generate)
}
func (c *PipelineRuleAgentController) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("POST", "/data/pipeline-rules/generate-from-excel", "data:service:rule:add")
}
func (c *PipelineRuleAgentController) Generate(ctx *gin.Context) {
	header, err := ctx.FormFile("file")
	if err != nil || header == nil {
		baizeContext.ParameterError(ctx)
		return
	}
	if header.Size > 20*1024*1024 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "Excel 文件不能超过 20MB"})
		return
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".xlsx" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "当前仅支持 .xlsx 文件；请先将旧版 .xls 另存为 .xlsx"})
		return
	}
	config, err := c.agent.Generate(ctx, header, "file")
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	config, err = c.rules.ValidateConfig(ctx, config, "file")
	if err != nil {
		baizeContext.Waring(ctx, "Agent 生成的 Pipeline 配置校验失败: "+err.Error())
		return
	}
	baizeContext.SuccessData(ctx, gin.H{"config": config})
}
