package routes

import (
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

func NewMcpRouter(r *gin.Engine) *gin_mcp.GinMCP {
	// --- Configure MCP Server ---
	mpcServer := gin_mcp.New(r, &gin_mcp.Config{
		Name:        "Product API",
		Description: "API for managing products.",
		BaseURL:     "http://localhost:8080",
	})
	//type McpConfig struct {
	//	Path           string `mapstructure:"path"`
	//	OperationsPath string `mapstructure:"operationsPath"`
	//}
	//// 把读取到的配置信息反序列化到 Conf 变量中
	//var mcpConfig McpConfig
	//if err := viper.UnmarshalKey("mcp", &mcpConfig); err != nil {
	//	panic(err)
	//}
	//
	//// 4. Mount the MCP server endpoint
	//mpcServer.Mount(mcpConfig.Path,
	//	mcpConfig.OperationsPath) // MCP clients will connect here
	return mpcServer
}
