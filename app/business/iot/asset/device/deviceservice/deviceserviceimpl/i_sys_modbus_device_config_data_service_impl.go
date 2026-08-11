package deviceserviceimpl

import (
	"strconv"

	"go.uber.org/zap"
	"nova-factory-server/app/business/iot/asset/device/devicedao"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/asset/device/deviceservice"
	"nova-factory-server/app/business/iot/metric/device/metricdao"
	metricentity "nova-factory-server/app/business/iot/metric/device/metricmodels/entity"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
)

type ISysModbusDeviceConfigDataServiceImpl struct {
	dao       devicedao.ISysModbusDeviceConfigDataDao
	metricDao metricdao.IMetricDao
}

func NewISysModbusDeviceConfigDataServiceImpl(dao devicedao.ISysModbusDeviceConfigDataDao, metricDao metricdao.IMetricDao) deviceservice.ISysModbusDeviceConfigDataService {
	return &ISysModbusDeviceConfigDataServiceImpl{
		dao:       dao,
		metricDao: metricDao,
	}
}
func (i *ISysModbusDeviceConfigDataServiceImpl) Add(c *gin.Context, template *devicemodels.SetSysModbusDeviceConfigDataReq) (*devicemodels.SysModbusDeviceConfigData, error) {
	data := devicemodels.OfSysModbusDeviceConfigData(template)
	data.DeviceConfigID = snowflake.GenID()
	data.DeptID = baizeContext.GetDeptId(c)
	data.SetCreateBy(baizeContext.GetUserId(c))
	value, err := i.dao.Add(c, data)
	if err != nil {
		return nil, err
	}
	if err := i.syncMetricTemplates(c, []uint64{uint64(data.TemplateID)}); err != nil {
		return nil, err
	}
	return value, nil
}
func (i *ISysModbusDeviceConfigDataServiceImpl) Update(c *gin.Context, template *devicemodels.SetSysModbusDeviceConfigDataReq) (*devicemodels.SysModbusDeviceConfigData, error) {
	data := devicemodels.OfSysModbusDeviceConfigData(template)
	data.SetUpdateBy(baizeContext.GetUserId(c))

	templateIds := []uint64{uint64(data.TemplateID)}
	if data.DeviceConfigID > 0 {
		oldData, err := i.dao.GetById(c, uint64(data.DeviceConfigID))
		if err != nil {
			zap.L().Error("get old device template data error", zap.Error(err))
			return nil, err
		}
		if oldData != nil && oldData.TemplateID > 0 {
			templateIds = append(templateIds, uint64(oldData.TemplateID))
		}
	}

	value, err := i.dao.Update(c, data)
	if err != nil {
		return nil, err
	}
	if err := i.syncMetricTemplates(c, templateIds); err != nil {
		return nil, err
	}
	return value, nil
}
func (i *ISysModbusDeviceConfigDataServiceImpl) Remove(c *gin.Context, ids []string) error {
	dataIds := make([]uint64, 0, len(ids))
	for _, id := range ids {
		dataId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return err
		}
		dataIds = append(dataIds, dataId)
	}

	list, err := i.dao.GetByIds(c, dataIds)
	if err != nil {
		return err
	}
	templateIds := make([]uint64, 0)
	for _, item := range list {
		if item != nil && item.TemplateID > 0 {
			templateIds = append(templateIds, uint64(item.TemplateID))
		}
	}

	if err := i.dao.Remove(c, ids); err != nil {
		return err
	}
	return i.syncMetricTemplates(c, templateIds)
}
func (i *ISysModbusDeviceConfigDataServiceImpl) List(c *gin.Context, req *devicemodels.SysModbusDeviceConfigDataListReq) (*devicemodels.SysModbusDeviceConfigDataListData, error) {
	return i.dao.List(c, req)
}

func (i *ISysModbusDeviceConfigDataServiceImpl) syncMetricTemplates(c *gin.Context, templateIds []uint64) error {
	if i.metricDao == nil || i.metricDao.Template() == nil {
		return nil
	}

	uniqueTemplateIds := make([]uint64, 0, len(templateIds))
	seen := make(map[uint64]struct{})
	for _, templateId := range templateIds {
		if templateId == 0 {
			continue
		}
		if _, ok := seen[templateId]; ok {
			continue
		}
		seen[templateId] = struct{}{}
		uniqueTemplateIds = append(uniqueTemplateIds, templateId)
	}
	if len(uniqueTemplateIds) == 0 {
		return nil
	}

	list, err := i.dao.GetByTemplateIds(c, uniqueTemplateIds)
	if err != nil {
		return err
	}
	metricTemplateData := devicemodels.FromTemplateDataToMetricTemplate(list)
	for _, templateId := range uniqueTemplateIds {
		if hasMetricTemplateData(metricTemplateData, int64(templateId)) {
			continue
		}
		metricTemplateData = append(metricTemplateData, metricentity.DeviceTemplate{TemplateId: int64(templateId)})
	}
	return i.metricDao.Template().Update(metricTemplateData)
}

func hasMetricTemplateData(data []metricentity.DeviceTemplate, templateId int64) bool {
	for _, item := range data {
		if item.TemplateId == templateId && item.DataId > 0 {
			return true
		}
	}
	return false
}
