package configurationController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/configuration/configurationModels"
	"nova-factory-server/app/business/configuration/configurationService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Configuration struct {
	service configurationService.ConfigurationService
}

func NewConfiguration(service configurationService.ConfigurationService) *Configuration {
	return &Configuration{
		service: service,
	}
}

func (c2 *Configuration) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/configuration/manager")
	group.GET("/list", middlewares.HasPermission("configuration:manager:list"), c2.List)
	group.POST("/set", middlewares.HasPermission("configuration:manager:set"), c2.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("configuration:manager:remove"), c2.Remove)
}

// List 组态列表
// @Summary 组态列表
// @Description 组态列表
// @Tags 组态/组态管理
// @Param  object query configurationModels.SysConfigurationReq true "组态列表参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "组态列表成功"
// @Router /configuration/manager/list [get]
func (c2 *Configuration) List(c *gin.Context) {
	req := new(configurationModels.SysConfigurationReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := c2.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// Set 保存组态
// @Summary 保存组态
// @Description 保存组态
// @Tags 组态/组态管理
// @Param  object body configurationModels.SetSysConfiguration true "设置组态参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置组态成功"
// @Router /configuration/manager/set [post]
func (c2 *Configuration) Set(c *gin.Context) {
	req := new(configurationModels.SetSysConfiguration)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	value, err := c2.service.Set(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, value)
}

// Remove 删除组态
// @Summary 删除组态
// @Description 删除组态
// @Tags 组态/组态管理
// @Param  ids path string true "ids"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "删除组态成功"
// @Router /configuration/manager/remove [delete]
func (c2 *Configuration) Remove(c *gin.Context) {
	recordIds := baizeContext.ParamStringArray(c, "ids")
	if len(recordIds) == 0 {
		baizeContext.Waring(c, "请选择要删除的组态")
		return
	}

	err := c2.service.Remove(c, recordIds)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}
	baizeContext.Success(c)
}
