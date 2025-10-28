package deviceServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
)

type IDeviceSubjectServiceImpl struct {
	dao deviceDao.IDeviceSubjectDao
}

func NewIDeviceSubjectServiceImpl(dao deviceDao.IDeviceSubjectDao) deviceService.IDeviceSubjectService {
	return &IDeviceSubjectServiceImpl{
		dao: dao,
	}
}

func (i *IDeviceSubjectServiceImpl) Set(c *gin.Context, data *deviceModels.SysDeviceSubjectVO) (*deviceModels.SysDeviceSubject, error) {
	return i.dao.Set(c, data)
}

func (i *IDeviceSubjectServiceImpl) List(c *gin.Context, req *deviceModels.SysDeviceSubjectReq) {
	i.dao.List(c, req)
}
func (i *IDeviceSubjectServiceImpl) Remove(c *gin.Context, ids []string) {
	i.dao.Remove(c, ids)
}
