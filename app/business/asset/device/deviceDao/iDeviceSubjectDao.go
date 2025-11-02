package deviceDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceModels"
)

type IDeviceSubjectDao interface {
	Set(c *gin.Context, data *deviceModels.SysDeviceSubjectVO) (*deviceModels.SysDeviceSubject, error)
	List(c *gin.Context, req *deviceModels.SysDeviceSubjectReq) (*deviceModels.SysDeviceSubjectListData, error)
	Remove(c *gin.Context, ids []string) error
	GetBySubjectCode(c *gin.Context, code string) (*deviceModels.SysDeviceSubject, error)
	GetBySubjectCodeByNotId(c *gin.Context, id int64, code string) (*deviceModels.SysDeviceSubject, error)
}
