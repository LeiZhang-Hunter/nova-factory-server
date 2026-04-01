package deviceservice

import (
	"nova-factory-server/app/business/iot/asset/device/devicemodels"

	"github.com/gin-gonic/gin"
)

type IDeviceTemplateService interface {
	Add(c *gin.Context, template *devicemodels.SysDeviceTemplateSetReq) (*devicemodels.SysDeviceTemplate, error)
	Update(c *gin.Context, template *devicemodels.SysDeviceTemplateSetReq) (*devicemodels.SysDeviceTemplate, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *devicemodels.SysDeviceTemplateDQL) (*devicemodels.SysDeviceTemplateListData, error)
	GetById(c *gin.Context, id int64) (*devicemodels.SysDeviceTemplate, error)
}
