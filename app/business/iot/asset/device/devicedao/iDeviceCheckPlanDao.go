package devicedao

import (
	"nova-factory-server/app/business/iot/asset/device/devicemodels"

	"github.com/gin-gonic/gin"
)

type IDeviceCheckPlanDao interface {
	Set(c *gin.Context, data *devicemodels.SysDeviceCheckPlanVO) (*devicemodels.SysDeviceCheckPlan, error)
	List(c *gin.Context, req *devicemodels.SysDeviceCheckPlanReq) (*devicemodels.SysDeviceCheckPlanList, error)
	Remove(c *gin.Context, ids []string) error
}
