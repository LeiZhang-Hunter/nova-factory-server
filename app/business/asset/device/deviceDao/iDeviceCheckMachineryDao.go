package deviceDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceModels"
)

type IDeviceCheckMachineryDao interface {
	Set(c *gin.Context, data *deviceModels.SysDeviceCheckMachineryVO) (*deviceModels.SysDeviceCheckMachinery, error)
	List(c *gin.Context, req *deviceModels.SysDeviceCheckMachineryReq) (*deviceModels.SysDeviceCheckMachineryList, error)
	Remove(c *gin.Context, ids []string) error
}
