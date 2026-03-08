package deviceService

import (
	"nova-factory-server/app/business/iot/asset/device/deviceModels"

	"github.com/gin-gonic/gin"
)

type IDeviceCheckSubjectService interface {
	Set(c *gin.Context, data *deviceModels.SysDeviceCheckSubjectVO) (*deviceModels.SysDeviceCheckSubject, error)
	List(c *gin.Context, req *deviceModels.SysDeviceCheckSubjectReq) (*deviceModels.SysDeviceCheckSubjectList, error)
	Remove(c *gin.Context, ids []string) error
}
