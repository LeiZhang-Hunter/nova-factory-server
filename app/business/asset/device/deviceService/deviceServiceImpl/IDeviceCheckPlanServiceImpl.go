package deviceServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
)

type IDeviceCheckPlanServiceImpl struct {
	dao deviceDao.IDeviceCheckPlanDao
}

func NewIDeviceCheckPlanServiceImpl(dao deviceDao.IDeviceCheckPlanDao) deviceService.IDeviceCheckPlanService {
	return &IDeviceCheckPlanServiceImpl{
		dao: dao,
	}
}

func (i *IDeviceCheckPlanServiceImpl) Set(c *gin.Context, data *deviceModels.SysDeviceCheckPlanVO) (*deviceModels.SysDeviceCheckPlan, error) {
	return i.dao.Set(c, data)
}
func (i *IDeviceCheckPlanServiceImpl) List(c *gin.Context, req *deviceModels.SysDeviceCheckPlanReq) (*deviceModels.SysDeviceCheckPlanList, error) {
	return i.dao.List(c, req)
}
func (i *IDeviceCheckPlanServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
