package deviceservice

import (
	"nova-factory-server/app/business/iot/asset/device/devicemodels"

	"github.com/gin-gonic/gin"
)

type IDeviceSubjectService interface {
	Set(c *gin.Context, data *devicemodels.SysDeviceSubjectVO) (*devicemodels.SysDeviceSubject, error)
	List(c *gin.Context, req *devicemodels.SysDeviceSubjectReq) (*devicemodels.SysDeviceSubjectListData, error)
	Remove(c *gin.Context, ids []string) error
	GetBySubjectCode(c *gin.Context, code string) (*devicemodels.SysDeviceSubject, error)
	GetBySubjectCodeByNotId(c *gin.Context, id int64, code string) (*devicemodels.SysDeviceSubject, error)
}
