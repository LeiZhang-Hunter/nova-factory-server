package deviceserviceimpl

import (
	"errors"
	"go.uber.org/zap"
	deviceDao2 "nova-factory-server/app/business/iot/asset/device/devicedao"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/asset/device/deviceservice"
	"nova-factory-server/app/business/iot/metric/device/metricdao"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IDeviceTemplateServiceImpl struct {
	dao       deviceDao2.IDeviceTemplateDao
	dataDao   deviceDao2.ISysModbusDeviceConfigDataDao
	metricDao metricdao.IMetricDao
	db        *gorm.DB
}

func NewDeviceTemplateServiceImpl(dao deviceDao2.IDeviceTemplateDao,
	dataDao deviceDao2.ISysModbusDeviceConfigDataDao,
	db *gorm.DB, metricDao metricdao.IMetricDao) deviceservice.IDeviceTemplateService {
	return &IDeviceTemplateServiceImpl{
		dao:       dao,
		dataDao:   dataDao,
		db:        db,
		metricDao: metricDao,
	}
}

// Add 添加设备模板
func (i *IDeviceTemplateServiceImpl) Add(c *gin.Context, template *devicemodels.SysDeviceTemplateSetReq) (*devicemodels.SysDeviceTemplate, error) {
	data := devicemodels.ToSysDeviceTemplate(template)
	data.TemplateID = snowflake.GenID()
	data.DeptID = baizeContext.GetDeptId(c)
	data.SetCreateBy(baizeContext.GetUserId(c))
	return i.dao.Add(c, data)
}

// Update 更新设备模板
func (i *IDeviceTemplateServiceImpl) Update(c *gin.Context, template *devicemodels.SysDeviceTemplateSetReq) (*devicemodels.SysDeviceTemplate, error) {
	data := devicemodels.ToSysDeviceTemplate(template)
	data.SetUpdateBy(baizeContext.GetUserId(c))

	var value *devicemodels.SysDeviceTemplate
	err := i.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		var err error
		value, err = i.dao.UpdateWithTx(c, tx, data)
		if err != nil {
			zap.L().Error("update device template error", zap.Error(err))
			return err
		}

		list, err := i.dataDao.GetByTemplateIdsWithTx(c, tx, []uint64{uint64(data.TemplateID)})
		if err != nil {
			zap.L().Error("get device template data error", zap.Error(err))
			return err
		}

		// 修改时序数据库设备模板
		if i.metricDao.Template() == nil {
			return nil
		}

		metricTemplateData := devicemodels.FromTemplateDataToMetricTemplate(list)
		return i.metricDao.Template().Update(metricTemplateData)
	})
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Remove 删除设备模板
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
func (i *IDeviceTemplateServiceImpl) List(c *gin.Context, req *devicemodels.SysDeviceTemplateDQL) (*devicemodels.SysDeviceTemplateListData, error) {
	return i.dao.List(c, req)
}

func (i *IDeviceTemplateServiceImpl) GetById(c *gin.Context, id int64) (*devicemodels.SysDeviceTemplate, error) {
	return i.dao.GetById(c, id)
}
