package devicemonitordao

import (
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"

	"github.com/gin-gonic/gin"
)

type IDeviceDataReportDao interface {
	DevList(c *gin.Context) ([]devicemonitormodel.SysIotDbDevMapData, error)
	GetDevList(c *gin.Context, dev []string) ([]devicemonitormodel.SysIotDbDevMapData, error)
	Save(c *gin.Context, data *devicemonitormodel.SysIotDbDevMap) error
	Remove(c *gin.Context, dev string) error
	List(c *gin.Context, req *devicemonitormodel.DevListReq) (*devicemonitormodel.DevListResp, error)
}
