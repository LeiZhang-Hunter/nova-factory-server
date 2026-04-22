package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"nova-factory-server/app/utils/gin_mcp"
)

type generatorConfig struct {
	Name              string `mapstructure:"name"`
	Description       string `mapstructure:"description"`
	BaseURL           string `mapstructure:"baseURL"`
	ScanPath          string `mapstructure:"scanPath"`
	ToolsOutputPath   string `mapstructure:"toolsOutputPath"`
	OperationsOutPath string `mapstructure:"operationsOutputPath"`
}

func main() {
	cfg, err := loadGeneratorConfig()
	if err != nil {
		panic(err)
	}

	// --- Configure MCP Server ---
	r := gin.New()
	mpcServer := gin_mcp.New(r, &gin_mcp.Config{
		Name:        cfg.Name,
		Description: cfg.Description,
		BaseURL:     cfg.BaseURL,
		Path:        cfg.ScanPath,
	})
	err = mpcServer.SetupServer()
	if err != nil {
		panic(err)
	}
	tools := mpcServer.GetTools()
	toolsContent, err := json.MarshalIndent(tools, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(tools)
	if err = writeFileWithMkdir(cfg.ToolsOutputPath, toolsContent); err != nil {
		panic(err)
	}

	operationsContent, err := json.MarshalIndent(mpcServer.GetOperations(), "", "\t")
	if err != nil {
		panic(err)
	}
	if err = writeFileWithMkdir(cfg.OperationsOutPath, operationsContent); err != nil {
		panic(err)
	}
}

func loadGeneratorConfig() (*generatorConfig, error) {
	configPath := pflag.String("config", "./config/config.yaml", "配置文件路径")
	name := pflag.String("name", "", "MCP server 名称")
	description := pflag.String("description", "", "MCP server 描述")
	baseURL := pflag.String("base-url", "", "MCP server BaseURL")
	scanPath := pflag.String("scan-path", "", "GinMCP 扫描路径")
	toolsOutputPath := pflag.String("tools-output", "", "生成 mcp.json 输出路径")
	operationsOutputPath := pflag.String("operations-output", "", "生成 operations.json 输出路径")
	pflag.Parse()

	viper.SetConfigFile(*configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &generatorConfig{
		Name:        "Product API",
		Description: "API for managing products.",
		BaseURL:     "http://localhost:8080",
		ScanPath:    "../app",
	}
	_ = viper.UnmarshalKey("mcpGen", cfg)

	if cfg.ToolsOutputPath == "" {
		cfg.ToolsOutputPath = viper.GetString("mcp.path")
	}
	if cfg.OperationsOutPath == "" {
		cfg.OperationsOutPath = viper.GetString("mcp.operationsPath")
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = viper.GetString("host")
	}

	overrideString(&cfg.Name, *name)
	overrideString(&cfg.Description, *description)
	overrideString(&cfg.BaseURL, *baseURL)
	overrideString(&cfg.ScanPath, *scanPath)
	overrideString(&cfg.ToolsOutputPath, *toolsOutputPath)
	overrideString(&cfg.OperationsOutPath, *operationsOutputPath)

	if cfg.ScanPath == "" {
		cfg.ScanPath = "../app"
	}
	if cfg.ToolsOutputPath == "" {
		cfg.ToolsOutputPath = "./config/mcp.json"
	}
	if cfg.OperationsOutPath == "" {
		cfg.OperationsOutPath = "./config/operations.json"
	}
	return cfg, nil
}

func overrideString(target *string, value string) {
	if value == "" {
		return
	}
	*target = value
}

func writeFileWithMkdir(outputPath string, content []byte) error {
	dir := filepath.Dir(outputPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(outputPath, content, 0644)
}
