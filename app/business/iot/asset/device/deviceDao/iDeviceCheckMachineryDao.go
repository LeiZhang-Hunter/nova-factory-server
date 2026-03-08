package deviceDao

import (
	"nova-factory-server/app/business/iot/asset/device/deviceModels"

	"github.com/gin-gonic/gin"
)

type IDeviceCheckMachineryDao interface {
	Set(c *gin.Context, data *deviceModels.SysDeviceCheckMachineryVO) (*deviceModels.SysDeviceCheckMachinery, error)
	List(c *gin.Context, req *deviceModels.SysDeviceCheckMachineryReq) (*deviceModels.SysDeviceCheckMachineryList, error)
	Remove(c *gin.Context, ids []string) error
}
