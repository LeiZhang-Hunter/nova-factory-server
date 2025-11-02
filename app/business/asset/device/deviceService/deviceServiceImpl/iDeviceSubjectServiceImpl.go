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

func (i *IDeviceSubjectServiceImpl) List(c *gin.Context, req *deviceModels.SysDeviceSubjectReq) (*deviceModels.SysDeviceSubjectListData, error) {
	return i.dao.List(c, req)
}

func (i *IDeviceSubjectServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}

func (i *IDeviceSubjectServiceImpl) GetBySubjectCode(c *gin.Context, code string) (*deviceModels.SysDeviceSubject, error) {
	return i.dao.GetBySubjectCode(c, code)
}

func (i *IDeviceSubjectServiceImpl) GetBySubjectCodeByNotId(c *gin.Context, id int64, code string) (*deviceModels.SysDeviceSubject, error) {
	return i.dao.GetBySubjectCodeByNotId(c, id, code)
}
