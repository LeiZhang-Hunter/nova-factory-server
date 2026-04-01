package devicemonitorservice

import (
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"

	"github.com/gin-gonic/gin"
)

type IDeviceDataReportService interface {
	DevList(c *gin.Context) ([]devicemonitormodel.SysIotDbDevMapData, error)
	GetDevList(c *gin.Context, req *devicemonitormodel.DevListReq) (*devicemonitormodel.DevListResp, error)
	// SetDevMap 设置测点
	SetDevMap(c *gin.Context, info *devicemonitormodel.SetDevMapInfo) error
	// RemoveDevMap 删除测点
	RemoveDevMap(c *gin.Context, dev string) error
}
