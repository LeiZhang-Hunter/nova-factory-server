package daemonizeController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
	"nova-factory-server/app/business/daemonize/daemonizeService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gateway/v1/config/pipeline"
	"nova-factory-server/app/utils/yaml"
)

type Config struct {
	service       daemonizeService.IGatewayConfigService
	configService daemonizeService.IotAgentConfigService
	agentService  daemonizeService.IotAgentService
}

func NewConfig(service daemonizeService.IGatewayConfigService,
	configService daemonizeService.IotAgentConfigService,
	agentService daemonizeService.IotAgentService) *Config {
	return &Config{
		service:       service,
		configService: configService,
		agentService:  agentService,
	}
}

func (c *Config) PrivateRoutes(router *gin.RouterGroup) {
	routers := router.Group("/gateway/agent/config")
	routers.GET("/generate", middlewares.HasPermission("gateway:agent:config:generate"), c.Generate) // 生成配置
	routers.GET("/list", middlewares.HasPermission("gateway:agent:config:list"), c.List)             // 配置列表
	routers.POST("/set", middlewares.HasPermission("gateway:agent:config:set"), c.Set)               // 保存配置
}

func (s *Config) PublicRoutes(router *gin.RouterGroup) {
	routers := router.Group("/api/gateway/agent/config/v1")
	routers.GET("/info", middlewares.HasPermission("gateway:agent:config:list"), s.List)
}

// Generate 生成Agent配置
// @Summary 生成Agent配置
// @Description 生成Agent配置
// @Tags 网关管理/Agent管理
// @Param  object query daemonizeModels.GenerateGatewayConfigReq true "参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /gateway/agent/config/generate [get]
func (c *Config) Generate(ctx *gin.Context) {
	req := new(daemonizeModels.GenerateGatewayConfigReq)
	err := ctx.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(ctx)
		return
	}
	config, err := c.service.Generate(ctx, int64(req.ObjectID))
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	content, err := yaml.Marshal(config)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	var res daemonizeModels.GenerateGatewayConfigRes
	res.Data = string(content)
	baizeContext.SuccessData(ctx, res)
	return
}

// List Agent配置列表
// @Summary Agent配置列表
// @Description Agent配置列表
// @Tags 网关管理/Agent管理
// @Param  object query daemonizeModels.SysIotAgentConfigListReq true "参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /gateway/agent/config/list [get]
func (c *Config) List(ctx *gin.Context) {
	req := new(daemonizeModels.SysIotAgentConfigListReq)
	err := ctx.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(ctx)
		return
	}
	list, err := c.configService.List(ctx, req)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, list)
}

// Set 设置Agent配置
// @Summary 设置Agent配置
// @Description 设置Agent配置
// @Tags 网关管理/Agent管理
// @Param  object body daemonizeModels.SysIotAgentConfigSetReq true "参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /gateway/agent/config/set [post]
func (c *Config) Set(ctx *gin.Context) {
	req := new(daemonizeModels.SysIotAgentConfigSetReq)
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(ctx)
		return
	}

	var pipelines pipeline.PipelineConfig
	err = yaml.Unmarshal([]byte(req.Content), &pipelines)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}

	info, err := c.agentService.GetByObjectId(ctx, uint64(req.AgentObjectID))
	if err != nil {
		baizeContext.Waring(ctx, "配置保存失败")
		return
	}
	if info == nil {
		baizeContext.Waring(ctx, "agent不存在")
		return
	}
	if req.ID == 0 {
		data, err := c.configService.Create(ctx, req)
		if err != nil {
			baizeContext.Waring(ctx, "配置保存失败")
			return
		}

		err = c.agentService.UpdateLastConfig(ctx, data.ID, []uint64{uint64(req.AgentObjectID)})
		if err != nil {
			zap.L().Error("UpdateLastConfig error", zap.Error(err))
		}

		// 更新gateway的最后一个配置id
		baizeContext.SuccessData(ctx, data)
	} else {
		data, err := c.configService.Update(ctx, req)
		if err != nil {
			baizeContext.Waring(ctx, "配置保存失败")
			return
		}
		baizeContext.SuccessData(ctx, data)
	}
}
