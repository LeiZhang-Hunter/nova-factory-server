package deviceserviceimpl

import (
	"nova-factory-server/app/business/iot/asset/device/devicedao"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/asset/device/deviceservice"

	"github.com/gin-gonic/gin"
)

type IDeviceSubjectServiceImpl struct {
	dao devicedao.IDeviceSubjectDao
}

func NewIDeviceSubjectServiceImpl(dao devicedao.IDeviceSubjectDao) deviceservice.IDeviceSubjectService {
	return &IDeviceSubjectServiceImpl{
		dao: dao,
	}
}

func (i *IDeviceSubjectServiceImpl) Set(c *gin.Context, data *devicemodels.SysDeviceSubjectVO) (*devicemodels.SysDeviceSubject, error) {
	return i.dao.Set(c, data)
}

func (i *IDeviceSubjectServiceImpl) List(c *gin.Context, req *devicemodels.SysDeviceSubjectReq) (*devicemodels.SysDeviceSubjectListData, error) {
	return i.dao.List(c, req)
}

func (i *IDeviceSubjectServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}

func (i *IDeviceSubjectServiceImpl) GetBySubjectCode(c *gin.Context, code string) (*devicemodels.SysDeviceSubject, error) {
	return i.dao.GetBySubjectCode(c, code)
}

func (i *IDeviceSubjectServiceImpl) GetBySubjectCodeByNotId(c *gin.Context, id int64, code string) (*devicemodels.SysDeviceSubject, error) {
	return i.dao.GetBySubjectCodeByNotId(c, id, code)
}
