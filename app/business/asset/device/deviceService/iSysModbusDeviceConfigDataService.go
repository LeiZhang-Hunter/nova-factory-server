package deviceService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceModels"
)

type ISysModbusDeviceConfigDataService interface {
	Add(c *gin.Context, template *deviceModels.SetSysModbusDeviceConfigDataReq) (*deviceModels.SysModbusDeviceConfigData, error)
	Update(c *gin.Context, template *deviceModels.SetSysModbusDeviceConfigDataReq) (*deviceModels.SysModbusDeviceConfigData, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *deviceModels.SysModbusDeviceConfigDataListReq) (*deviceModels.SysModbusDeviceConfigDataListData, error)
}
