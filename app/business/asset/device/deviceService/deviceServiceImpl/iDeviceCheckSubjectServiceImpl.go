package deviceServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
)

type IDeviceCheckSubjectServiceImpl struct {
	dao deviceDao.IDeviceCheckSubjectDao
}

func NewIDeviceCheckSubjectServiceImpl(dao deviceDao.IDeviceCheckSubjectDao) deviceService.IDeviceCheckSubjectService {
	return &IDeviceCheckSubjectServiceImpl{
		dao: dao,
	}
}

func (i *IDeviceCheckSubjectServiceImpl) Set(c *gin.Context, data *deviceModels.SysDeviceCheckSubjectVO) (*deviceModels.SysDeviceCheckSubject, error) {
	return i.dao.Set(c, data)
}
func (i *IDeviceCheckSubjectServiceImpl) List(c *gin.Context, req *deviceModels.SysDeviceCheckSubjectReq) (*deviceModels.SysDeviceCheckSubjectList, error) {
	return i.dao.List(c, req)
}

func (i *IDeviceCheckSubjectServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
