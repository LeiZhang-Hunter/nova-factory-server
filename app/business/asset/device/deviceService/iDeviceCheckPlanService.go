package deviceService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceModels"
)

type IDeviceCheckPlanService interface {
	Set(c *gin.Context, data *deviceModels.SysDeviceCheckPlanVO) (*deviceModels.SysDeviceCheckPlan, error)
	List(c *gin.Context, req *deviceModels.SysDeviceCheckPlanReq) (*deviceModels.SysDeviceCheckPlanList, error)
	Remove(c *gin.Context, ids []string) error
}
