package routes

import (
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

var ProviderSet = wire.NewSet(NewGinEngine)

func NewGinEngine() *gin.Engine {

	if setting.Conf.Mode != "dev" {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
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

	// pprof
	pprof.Register(r)

	//product.Laboratory.PrivateMcpRoutes(mpcServer)

	return r

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
