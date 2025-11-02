package deviceServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
)

type IDeviceCheckPlanServiceImpl struct {
	service deviceService.IDeviceCheckPlanService
}

func NewIDeviceCheckPlanServiceImpl(service deviceService.IDeviceCheckPlanService) deviceService.IDeviceCheckPlanService {
	return &IDeviceCheckPlanServiceImpl{
		service: service,
	}
}

func (i *IDeviceCheckPlanServiceImpl) Set(c *gin.Context, data *deviceModels.SysDeviceCheckPlanVO) (*deviceModels.SysDeviceCheckPlan, error) {
	return i.service.Set(c, data)
}
func (i *IDeviceCheckPlanServiceImpl) List(c *gin.Context, req *deviceModels.SysDeviceCheckPlanReq) (*deviceModels.SysDeviceCheckPlanList, error) {
	return i.service.List(c, req)
}
func (i *IDeviceCheckPlanServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.service.Remove(c, ids)
}
