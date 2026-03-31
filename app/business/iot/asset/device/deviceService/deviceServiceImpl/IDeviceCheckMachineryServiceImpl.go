package deviceServiceImpl

import (
	"nova-factory-server/app/business/iot/asset/device/deviceDao"
	"nova-factory-server/app/business/iot/asset/device/deviceModels"
	"nova-factory-server/app/business/iot/asset/device/deviceService"

	"github.com/gin-gonic/gin"
)

type IDeviceCheckMachineryServiceImpl struct {
	dao deviceDao.IDeviceCheckMachineryDao
}

func NewIDeviceCheckMachineryServiceImpl(dao deviceDao.IDeviceCheckMachineryDao) deviceService.IDeviceCheckMachineryService {
	return &IDeviceCheckMachineryServiceImpl{
		dao: dao,
	}
}

func (i *IDeviceCheckMachineryServiceImpl) Set(c *gin.Context, data *deviceModels.SysDeviceCheckMachineryVO) (*deviceModels.SysDeviceCheckMachinery, error) {
	return i.dao.Set(c, data)
}
func (i *IDeviceCheckMachineryServiceImpl) List(c *gin.Context, req *deviceModels.SysDeviceCheckMachineryReq) (*deviceModels.SysDeviceCheckMachineryList, error) {
	return i.dao.List(c, req)
}
func (i *IDeviceCheckMachineryServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
