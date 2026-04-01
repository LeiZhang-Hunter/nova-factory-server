package deviceserviceimpl

import (
	"nova-factory-server/app/business/iot/asset/device/devicedao"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/asset/device/deviceservice"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
)

type ISysModbusDeviceConfigDataServiceImpl struct {
	dao devicedao.ISysModbusDeviceConfigDataDao
}

func NewISysModbusDeviceConfigDataServiceImpl(dao devicedao.ISysModbusDeviceConfigDataDao) deviceservice.ISysModbusDeviceConfigDataService {
	return &ISysModbusDeviceConfigDataServiceImpl{
		dao: dao,
	}
}
func (i *ISysModbusDeviceConfigDataServiceImpl) Add(c *gin.Context, template *devicemodels.SetSysModbusDeviceConfigDataReq) (*devicemodels.SysModbusDeviceConfigData, error) {
	data := devicemodels.OfSysModbusDeviceConfigData(template)
	data.DeviceConfigID = snowflake.GenID()
	data.DeptID = baizeContext.GetDeptId(c)
	data.SetCreateBy(baizeContext.GetUserId(c))
	return i.dao.Add(c, data)
}
func (i *ISysModbusDeviceConfigDataServiceImpl) Update(c *gin.Context, template *devicemodels.SetSysModbusDeviceConfigDataReq) (*devicemodels.SysModbusDeviceConfigData, error) {
	data := devicemodels.OfSysModbusDeviceConfigData(template)
	data.SetUpdateBy(baizeContext.GetUserId(c))
	return i.dao.Update(c, data)
}
func (i *ISysModbusDeviceConfigDataServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
func (i *ISysModbusDeviceConfigDataServiceImpl) List(c *gin.Context, req *devicemodels.SysModbusDeviceConfigDataListReq) (*devicemodels.SysModbusDeviceConfigDataListData, error) {
	return i.dao.List(c, req)
}
