package devicedao

import (
	"nova-factory-server/app/business/iot/asset/device/devicemodels"

	"github.com/gin-gonic/gin"
)

type ISysModbusDeviceConfigDataDao interface {
	Add(c *gin.Context, data *devicemodels.SysModbusDeviceConfigData) (*devicemodels.SysModbusDeviceConfigData, error)
	Update(c *gin.Context, data *devicemodels.SysModbusDeviceConfigData) (*devicemodels.SysModbusDeviceConfigData, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *devicemodels.SysModbusDeviceConfigDataListReq) (*devicemodels.SysModbusDeviceConfigDataListData, error)
	GetByTemplateIds(c *gin.Context, ids []uint64) ([]*devicemodels.SysModbusDeviceConfigData, error)
	GetById(c *gin.Context, id uint64) (*devicemodels.SysModbusDeviceConfigData, error)
	GetByIds(c *gin.Context, id []uint64) ([]*devicemodels.SysModbusDeviceConfigData, error)
}
