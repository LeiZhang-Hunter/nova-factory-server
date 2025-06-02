package deviceService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceModels"
)

type IDeviceTemplateService interface {
	Add(c *gin.Context, template *deviceModels.SysDeviceTemplateSetReq) (*deviceModels.SysDeviceTemplate, error)
	Update(c *gin.Context, template *deviceModels.SysDeviceTemplateSetReq) (*deviceModels.SysDeviceTemplate, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *deviceModels.SysDeviceTemplateDQL) (*deviceModels.SysDeviceTemplateListData, error)
}
