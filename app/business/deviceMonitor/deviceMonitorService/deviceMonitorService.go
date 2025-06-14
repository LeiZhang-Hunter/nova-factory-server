package deviceMonitorService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceModels"
)

type DeviceMonitorService interface {
	List(c *gin.Context) (*deviceModels.DeviceInfoListData, error)
}
