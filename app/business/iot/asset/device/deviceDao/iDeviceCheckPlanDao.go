package deviceDao

import (
	"nova-factory-server/app/business/iot/asset/device/deviceModels"

	"github.com/gin-gonic/gin"
)

type IDeviceCheckPlanDao interface {
	Set(c *gin.Context, data *deviceModels.SysDeviceCheckPlanVO) (*deviceModels.SysDeviceCheckPlan, error)
	List(c *gin.Context, req *deviceModels.SysDeviceCheckPlanReq) (*deviceModels.SysDeviceCheckPlanList, error)
	Remove(c *gin.Context, ids []string) error
}
