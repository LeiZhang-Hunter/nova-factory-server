package deviceServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type IDeviceTemplateServiceImpl struct {
	dao deviceDao.IDeviceTemplateDao
}

func NewDeviceTemplateServiceImpl(dao deviceDao.IDeviceTemplateDao) deviceService.IDeviceTemplateService {
	return &IDeviceTemplateServiceImpl{
		dao: dao,
	}
}
func (i *IDeviceTemplateServiceImpl) Add(c *gin.Context, template *deviceModels.SysDeviceTemplateSetReq) (*deviceModels.SysDeviceTemplate, error) {
	data := deviceModels.ToSysDeviceTemplate(template)
	data.TemplateID = snowflake.GenID()
	data.DeptID = baizeContext.GetDeptId(c)
	data.SetCreateBy(baizeContext.GetUserId(c))
	return i.dao.Add(c, data)
}
func (i *IDeviceTemplateServiceImpl) Update(c *gin.Context, template *deviceModels.SysDeviceTemplateSetReq) (*deviceModels.SysDeviceTemplate, error) {
	data := deviceModels.ToSysDeviceTemplate(template)
	data.SetUpdateBy(baizeContext.GetUserId(c))
	return i.dao.Update(c, data)
}
func (i *IDeviceTemplateServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
func (i *IDeviceTemplateServiceImpl) List(c *gin.Context, req *deviceModels.SysDeviceTemplateDQL) (*deviceModels.SysDeviceTemplateListData, error) {
	return i.dao.List(c, req)
}
