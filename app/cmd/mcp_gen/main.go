package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/utils/gin_mcp"
	"os"
)

func main() {
	// --- Configure MCP Server ---
	r := gin.New()
	mpcServer := gin_mcp.New(r, &gin_mcp.Config{
		Name:        "Product API",
		Description: "API for managing products.",
		BaseURL:     "http://localhost:8080",
		Path:        "/home/zhanglei/project/zhanglei/nova-factory-server/app",
	})
	err := mpcServer.SetupServer()
	if err != nil {
		panic(err)
	}
	tools := mpcServer.GetTools()
	toolsContent, err := json.MarshalIndent(tools, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(tools)
	err = os.WriteFile("/home/zhanglei/project/zhanglei/nova-factory-server/config/mcp.json", []byte(toolsContent), 0644) // 0644是文件权限
	if err != nil {
		panic(err)
	}

	operationsContent, err := json.MarshalIndent(mpcServer.GetOperations(), "", "\t")
	err = os.WriteFile("/home/zhanglei/project/zhanglei/nova-factory-server/config/operations.json", []byte(operationsContent), 0644) // 0644是文件权限
	if err != nil {
		panic(err)
	}
}
