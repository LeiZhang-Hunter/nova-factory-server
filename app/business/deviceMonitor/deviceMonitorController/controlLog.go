package deviceMonitorController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type ControlLog struct {
	service deviceMonitorService.ControlLogService
}

func NewControlLog(service deviceMonitorService.ControlLogService) *ControlLog {
	return &ControlLog{
		service: service,
	}
}

func (that *ControlLog) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/device/monitor/control/log")
	group.GET("/list", middlewares.HasPermission("metric:monitor:control:log:list"), that.List)
}

// List 控制记录列表
// @Summary 控制记录列表
// @Description 控制记录列表
// @Param  object query deviceMonitorModel.ControlLogListReq true "控制记录列表"
// @Tags 设备接入/控制记录
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /device/monitor/control/log/list [get]
func (that *ControlLog) List(c *gin.Context) {
	req := new(deviceMonitorModel.ControlLogListReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := that.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}
