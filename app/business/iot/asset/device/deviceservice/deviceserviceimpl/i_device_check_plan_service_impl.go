package deviceserviceimpl

import (
	"nova-factory-server/app/business/iot/asset/device/devicedao"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/asset/device/deviceservice"

	"github.com/gin-gonic/gin"
)

type IDeviceCheckPlanServiceImpl struct {
	dao devicedao.IDeviceCheckPlanDao
}

func NewIDeviceCheckPlanServiceImpl(dao devicedao.IDeviceCheckPlanDao) deviceservice.IDeviceCheckPlanService {
	return &IDeviceCheckPlanServiceImpl{
		dao: dao,
	}
}

func (i *IDeviceCheckPlanServiceImpl) Set(c *gin.Context, data *devicemodels.SysDeviceCheckPlanVO) (*devicemodels.SysDeviceCheckPlan, error) {
	return i.dao.Set(c, data)
}
func (i *IDeviceCheckPlanServiceImpl) List(c *gin.Context, req *devicemodels.SysDeviceCheckPlanReq) (*devicemodels.SysDeviceCheckPlanList, error) {
	return i.dao.List(c, req)
}
func (i *IDeviceCheckPlanServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
