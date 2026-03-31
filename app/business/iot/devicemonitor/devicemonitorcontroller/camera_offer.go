package devicemonitorcontroller

import (
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	client "github.com/novawatcher-io/nova-factory-payload/camera/v1"
)

func (d *DeviceMonitor) privateCameraOfferRoutes(monitor *gin.RouterGroup) {
	monitor.POST("/camera/offer", middlewares.HasPermission("device:monitor:control"), d.CameraOffer)
}

// CameraOffer 摄像头 WebRTC SDP 协商
// @Summary 摄像头 WebRTC SDP 协商
// @Description 提交前端 Offer，返回播放地址与 Answer
// @Tags 设备监控/摄像头
// @Param object body devicemonitormodel.CameraOfferReq true "摄像头协商参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "协商成功"
// @Router /device/monitor/camera/offer [post]
func (d *DeviceMonitor) CameraOffer(c *gin.Context) {
	req := new(devicemonitormodel.CameraOfferReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	timeout := 10 * time.Second
	if req.TimeoutMS > 0 {
		timeout = time.Duration(req.TimeoutMS) * time.Millisecond
	}

	ack, err := d.cameraGrpc.PublishStartByBroadcast(req.Node, &client.SubscribeMessage{
		DeviceId:  req.DeviceId,
		ChannelId: req.ChannelId,
		Sdp64:     req.SDP64,
	}, timeout)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	if strings.HasPrefix(ack.Token, "error:") {
		baizeContext.Waring(c, strings.TrimPrefix(ack.Token, "error:"))
		return
	}

	baizeContext.SuccessData(c, &devicemonitormodel.CameraOfferRes{
		Token:   ack.Token,
		PlayURL: ack.PlayUrl,
		WhepURL: ack.WhepUrl,
		SDP64:   ack.Sdp64,
	})
}
