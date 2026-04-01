package deviceserviceimpl

import (
	"nova-factory-server/app/business/iot/asset/device/devicedao"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/asset/device/deviceservice"

	"github.com/gin-gonic/gin"
)

type IDeviceCheckMachineryServiceImpl struct {
	dao devicedao.IDeviceCheckMachineryDao
}

func NewIDeviceCheckMachineryServiceImpl(dao devicedao.IDeviceCheckMachineryDao) deviceservice.IDeviceCheckMachineryService {
	return &IDeviceCheckMachineryServiceImpl{
		dao: dao,
	}
}

func (i *IDeviceCheckMachineryServiceImpl) Set(c *gin.Context, data *devicemodels.SysDeviceCheckMachineryVO) (*devicemodels.SysDeviceCheckMachinery, error) {
	return i.dao.Set(c, data)
}
func (i *IDeviceCheckMachineryServiceImpl) List(c *gin.Context, req *devicemodels.SysDeviceCheckMachineryReq) (*devicemodels.SysDeviceCheckMachineryList, error) {
	return i.dao.List(c, req)
}
func (i *IDeviceCheckMachineryServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
