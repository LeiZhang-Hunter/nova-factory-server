package deviceservice

import (
	"nova-factory-server/app/business/iot/asset/device/devicemodels"

	"github.com/gin-gonic/gin"
)

type IDeviceCheckSubjectService interface {
	Set(c *gin.Context, data *devicemodels.SysDeviceCheckSubjectVO) (*devicemodels.SysDeviceCheckSubject, error)
	List(c *gin.Context, req *devicemodels.SysDeviceCheckSubjectReq) (*devicemodels.SysDeviceCheckSubjectList, error)
	Remove(c *gin.Context, ids []string) error
}
