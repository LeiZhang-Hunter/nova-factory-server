package deviceserviceimpl

import (
	"nova-factory-server/app/business/iot/asset/device/devicedao"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/asset/device/deviceservice"

	"github.com/gin-gonic/gin"
)

type IDeviceCheckSubjectServiceImpl struct {
	dao devicedao.IDeviceCheckSubjectDao
}

func NewIDeviceCheckSubjectServiceImpl(dao devicedao.IDeviceCheckSubjectDao) deviceservice.IDeviceCheckSubjectService {
	return &IDeviceCheckSubjectServiceImpl{
		dao: dao,
	}
}

func (i *IDeviceCheckSubjectServiceImpl) Set(c *gin.Context, data *devicemodels.SysDeviceCheckSubjectVO) (*devicemodels.SysDeviceCheckSubject, error) {
	return i.dao.Set(c, data)
}
func (i *IDeviceCheckSubjectServiceImpl) List(c *gin.Context, req *devicemodels.SysDeviceCheckSubjectReq) (*devicemodels.SysDeviceCheckSubjectList, error) {
	return i.dao.List(c, req)
}

func (i *IDeviceCheckSubjectServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
