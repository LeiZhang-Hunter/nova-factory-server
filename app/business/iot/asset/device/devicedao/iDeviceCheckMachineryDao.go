package devicedao

import (
	"nova-factory-server/app/business/iot/asset/device/devicemodels"

	"github.com/gin-gonic/gin"
)

type IDeviceCheckMachineryDao interface {
	Set(c *gin.Context, data *devicemodels.SysDeviceCheckMachineryVO) (*devicemodels.SysDeviceCheckMachinery, error)
	List(c *gin.Context, req *devicemodels.SysDeviceCheckMachineryReq) (*devicemodels.SysDeviceCheckMachineryList, error)
	Remove(c *gin.Context, ids []string) error
}
