package devicedao

import (
	"nova-factory-server/app/business/iot/asset/device/devicemodels"

	"github.com/gin-gonic/gin"
)

type IDeviceTemplateDao interface {
	Add(c *gin.Context, template *devicemodels.SysDeviceTemplate) (*devicemodels.SysDeviceTemplate, error)
	Update(c *gin.Context, template *devicemodels.SysDeviceTemplate) (*devicemodels.SysDeviceTemplate, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *devicemodels.SysDeviceTemplateDQL) (*devicemodels.SysDeviceTemplateListData, error)
	GetById(c *gin.Context, id int64) (*devicemodels.SysDeviceTemplate, error)
	GetByIds(c *gin.Context, ids []uint64) ([]*devicemodels.SysDeviceTemplate, error)
}
