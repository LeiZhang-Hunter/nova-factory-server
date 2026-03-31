package deviceservice

import (
	"nova-factory-server/app/business/iot/asset/device/devicemodels"

	"github.com/gin-gonic/gin"
)

type ISysModbusDeviceConfigDataService interface {
	Add(c *gin.Context, template *devicemodels.SetSysModbusDeviceConfigDataReq) (*devicemodels.SysModbusDeviceConfigData, error)
	Update(c *gin.Context, template *devicemodels.SetSysModbusDeviceConfigDataReq) (*devicemodels.SysModbusDeviceConfigData, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *devicemodels.SysModbusDeviceConfigDataListReq) (*devicemodels.SysModbusDeviceConfigDataListData, error)
}
