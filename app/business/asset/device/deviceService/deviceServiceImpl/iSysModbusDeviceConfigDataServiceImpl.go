package deviceServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type ISysModbusDeviceConfigDataServiceImpl struct {
	dao deviceDao.ISysModbusDeviceConfigDataDao
}

func NewISysModbusDeviceConfigDataServiceImpl(dao deviceDao.ISysModbusDeviceConfigDataDao) deviceService.ISysModbusDeviceConfigDataService {
	return &ISysModbusDeviceConfigDataServiceImpl{
		dao: dao,
	}
}
func (i *ISysModbusDeviceConfigDataServiceImpl) Add(c *gin.Context, template *deviceModels.SetSysModbusDeviceConfigDataReq) (*deviceModels.SysModbusDeviceConfigData, error) {
	data := deviceModels.OfSysModbusDeviceConfigData(template)
	data.DeviceConfigID = snowflake.GenID()
	data.DeptID = baizeContext.GetDeptId(c)
	data.SetCreateBy(baizeContext.GetUserId(c))
	return i.dao.Add(c, data)
}
func (i *ISysModbusDeviceConfigDataServiceImpl) Update(c *gin.Context, template *deviceModels.SetSysModbusDeviceConfigDataReq) (*deviceModels.SysModbusDeviceConfigData, error) {
	data := deviceModels.OfSysModbusDeviceConfigData(template)
	data.SetUpdateBy(baizeContext.GetUserId(c))
	return i.dao.Update(c, data)
}
func (i *ISysModbusDeviceConfigDataServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
func (i *ISysModbusDeviceConfigDataServiceImpl) List(c *gin.Context, req *deviceModels.SysModbusDeviceConfigDataListReq) (*deviceModels.SysModbusDeviceConfigDataListData, error) {
	return i.dao.List(c, req)
}
