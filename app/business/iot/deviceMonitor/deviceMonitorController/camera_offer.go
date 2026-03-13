package deviceMonitorController

import (
	"nova-factory-server/app/business/iot/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"time"

	"github.com/gin-gonic/gin"
)

func (d *DeviceMonitor) privateCameraOfferRoutes(monitor *gin.RouterGroup) {
	monitor.POST("/camera/offer", middlewares.HasPermission("device:monitor:control"), d.CameraOffer)
}

func (d *DeviceMonitor) CameraOffer(c *gin.Context) {
	req := new(deviceMonitorModel.CameraOfferReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	timeout := 10 * time.Second
	if req.TimeoutMS > 0 {
		timeout = time.Duration(req.TimeoutMS) * time.Millisecond
	}

	ack, err := d.cameraGrpc.PublishStart(req.Node, subscribeMessage{
		DeviceId:  req.DeviceId,
		ChannelId: req.ChannelId,
		SDP64:     req.SDP64,
	}, timeout)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}

	baizeContext.SuccessData(c, &deviceMonitorModel.CameraOfferRes{
		Token:   ack.Token,
		PlayURL: ack.PlayURL,
		WhepURL: ack.WhepURL,
		SDP64:   ack.SDP64,
	})
}
