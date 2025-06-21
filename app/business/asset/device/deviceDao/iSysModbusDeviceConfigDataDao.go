package deviceDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceModels"
)

type ISysModbusDeviceConfigDataDao interface {
	Add(c *gin.Context, data *deviceModels.SysModbusDeviceConfigData) (*deviceModels.SysModbusDeviceConfigData, error)
	Update(c *gin.Context, data *deviceModels.SysModbusDeviceConfigData) (*deviceModels.SysModbusDeviceConfigData, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *deviceModels.SysModbusDeviceConfigDataListReq) (*deviceModels.SysModbusDeviceConfigDataListData, error)
	GetByTemplateIds(c *gin.Context, ids []uint64) ([]*deviceModels.SysModbusDeviceConfigData, error)
	GetById(c *gin.Context, id uint64) (*deviceModels.SysModbusDeviceConfigData, error)
	GetByIds(c *gin.Context, id []uint64) ([]*deviceModels.SysModbusDeviceConfigData, error)
}
