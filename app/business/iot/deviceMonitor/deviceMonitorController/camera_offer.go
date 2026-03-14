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

// CameraOffer 摄像头 WebRTC SDP 协商
// @Summary 摄像头 WebRTC SDP 协商
// @Description 提交前端 Offer，返回播放地址与 Answer
// @Tags 设备监控/摄像头
// @Param object body deviceMonitorModel.CameraOfferReq true "摄像头协商参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "协商成功"
// @Router /device/monitor/camera/offer [post]
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
