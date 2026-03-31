package deviceService

import (
	"nova-factory-server/app/business/iot/asset/device/deviceModels"

	"github.com/gin-gonic/gin"
)

type IDeviceTemplateService interface {
	Add(c *gin.Context, template *deviceModels.SysDeviceTemplateSetReq) (*deviceModels.SysDeviceTemplate, error)
	Update(c *gin.Context, template *deviceModels.SysDeviceTemplateSetReq) (*deviceModels.SysDeviceTemplate, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *deviceModels.SysDeviceTemplateDQL) (*deviceModels.SysDeviceTemplateListData, error)
	GetById(c *gin.Context, id int64) (*deviceModels.SysDeviceTemplate, error)
}
