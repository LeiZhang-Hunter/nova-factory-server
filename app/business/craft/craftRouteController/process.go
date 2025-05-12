package craftRouteController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/middlewares"
)

type Process struct {
}

func NewProcess(craftService craftRouteService.ICraftRouteService) *Process {
	return &Process{}
}

func (p *Process) PrivateRoutes(router *gin.RouterGroup) {
	routers := router.Group("/craft/process")
	routers.GET("/list", middlewares.HasPermission("craft:process"), p.GetProcessList)                        // 知识库列表
	routers.POST("/set", middlewares.HasPermission("craft:process:set"), p.SetProcess)                        // 设置工艺路线
	routers.DELETE("/remove/:dataset_id", middlewares.HasPermission("craft:process:remove"), p.RemoveProcess) //移除工艺路线
}

func (p *Process) GetProcessList(c *gin.Context) {

}

// SetProcess 设置工序
// @Summary 设置工序
// @Description 设置工序2
// @Tags 工艺管理/工序管理
// @Param  object body craftRouteModels.SysProProcess true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/process/set [post]
func (p *Process) SetProcess(c *gin.Context) {

}

func (p *Process) RemoveProcess(c *gin.Context) {

}
