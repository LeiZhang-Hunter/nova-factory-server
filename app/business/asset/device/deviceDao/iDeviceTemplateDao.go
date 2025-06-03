package deviceDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceModels"
)

type IDeviceTemplateDao interface {
	Add(c *gin.Context, template *deviceModels.SysDeviceTemplate) (*deviceModels.SysDeviceTemplate, error)
	Update(c *gin.Context, template *deviceModels.SysDeviceTemplate) (*deviceModels.SysDeviceTemplate, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *deviceModels.SysDeviceTemplateDQL) (*deviceModels.SysDeviceTemplateListData, error)
	GetById(c *gin.Context, id int64) (*deviceModels.SysDeviceTemplate, error)
}
