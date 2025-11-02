package deviceDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceModels"
)

type IDeviceCheckSubjectDao interface {
	Set(c *gin.Context, data *deviceModels.SysDeviceCheckSubjectVO) (*deviceModels.SysDeviceCheckSubject, error)
	List(c *gin.Context, req *deviceModels.SysDeviceCheckSubjectReq) (*deviceModels.SysDeviceCheckSubjectList, error)
	Remove(c *gin.Context, ids []string) error
}
