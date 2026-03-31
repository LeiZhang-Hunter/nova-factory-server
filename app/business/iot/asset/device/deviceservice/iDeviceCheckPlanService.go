package deviceservice

import (
	"nova-factory-server/app/business/iot/asset/device/devicemodels"

	"github.com/gin-gonic/gin"
)

type IDeviceCheckPlanService interface {
	Set(c *gin.Context, data *devicemodels.SysDeviceCheckPlanVO) (*devicemodels.SysDeviceCheckPlan, error)
	List(c *gin.Context, req *devicemodels.SysDeviceCheckPlanReq) (*devicemodels.SysDeviceCheckPlanList, error)
	Remove(c *gin.Context, ids []string) error
}
