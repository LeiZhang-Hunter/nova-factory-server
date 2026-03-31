package routes

import (
	"nova-factory-server/app/utils/gin_mcp"
	"nova-factory-server/app/utils/logger"
	"time"

	"github.com/gin-contrib/pprof"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"

	"net/http"
	"nova-factory-server/app/docs"
	"nova-factory-server/app/setting"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/google/wire"

	"github.com/gin-gonic/gin"
)

type App struct {
	Engine    *gin.Engine
	McpServer *gin_mcp.GinMCP
}

var ProviderSet = wire.NewSet(NewGinApp)

func NewGinApp() *App {

	if setting.Conf.Mode != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.NewLoggerMiddlewareBuilder().
		IgnorePaths("/ping").Build())
	r.Use(newCors())
	group := r.Group("")
	if setting.Conf.Mode == "dev" {
		host := setting.Conf.Host
		docs.SwaggerInfo.Host = host[strings.Index(host, "//")+2:]
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
		group.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

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

	pprof.Register(r)
	return &App{
		Engine:    r,
		McpServer: mpcServer,
	}
}

func NewGinEngine(app *App) *gin.Engine {
	return app.Engine
}

func newCors() gin.HandlerFunc {
	ss := []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Cache-Control", "Content-Language", "Content-Type", "Expires", "Last-Modified", "Pragma", "FooBar"}
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} //允许访问的域名
	config.AllowMethods = []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	config.ExposeHeaders = ss
	config.MaxAge = time.Hour
	config.AllowCredentials = false
	return cors.New(config)
}
