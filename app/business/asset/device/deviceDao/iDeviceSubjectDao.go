package deviceDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceModels"
)

type IDeviceSubjectDao interface {
	Set(c *gin.Context, data *deviceModels.SysDeviceSubjectVO) (*deviceModels.SysDeviceSubject, error)
	List(c *gin.Context, req *deviceModels.SysDeviceSubjectReq)
	Remove(c *gin.Context, ids []string)
}
