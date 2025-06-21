package deviceServiceImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
	"strconv"
)

type IDeviceTemplateServiceImpl struct {
	dao     deviceDao.IDeviceTemplateDao
	dataDao deviceDao.ISysModbusDeviceConfigDataDao
}

func NewDeviceTemplateServiceImpl(dao deviceDao.IDeviceTemplateDao, dataDao deviceDao.ISysModbusDeviceConfigDataDao) deviceService.IDeviceTemplateService {
	return &IDeviceTemplateServiceImpl{
		dao:     dao,
		dataDao: dataDao,
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
	var templateIds []uint64 = make([]uint64, 0)
	for _, id := range ids {
		templateId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return err
		}
		templateIds = append(templateIds, templateId)
	}

	// 检查协议模板下面是否有数据，有数据不允许删除
	data, err := i.dataDao.GetByTemplateIds(c, templateIds)
	if err != nil {
		return err
	}

	if data == nil || len(data) > 0 {
		return errors.New("模板下面有数据，不能重置")
	}

	return i.dao.Remove(c, ids)
}
func (i *IDeviceTemplateServiceImpl) List(c *gin.Context, req *deviceModels.SysDeviceTemplateDQL) (*deviceModels.SysDeviceTemplateListData, error) {
	return i.dao.List(c, req)
}

func (i *IDeviceTemplateServiceImpl) GetById(c *gin.Context, id int64) (*deviceModels.SysDeviceTemplate, error) {
	return i.dao.GetById(c, id)
}
