package deviceMonitorServiceImpl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService"
	"nova-factory-server/app/business/metric/device/metricDao"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/constant/device"
	"nova-factory-server/app/datasource/cache"
)

type DeviceMonitorServiceImpl struct {
	dao                 deviceDao.IDeviceDao
	metricDao           metricDao.IMetricDao
	cache               cache.Cache
	deviceConfigDataDao deviceDao.ISysModbusDeviceConfigDataDao
}

func NewDeviceMonitorServiceImpl(dao deviceDao.IDeviceDao, cache cache.Cache, metricDao metricDao.IMetricDao, deviceConfigDataDao deviceDao.ISysModbusDeviceConfigDataDao) deviceMonitorService.DeviceMonitorService {
	return &DeviceMonitorServiceImpl{
		dao:                 dao,
		cache:               cache,
		metricDao:           metricDao,
		deviceConfigDataDao: deviceConfigDataDao,
	}
}

func (d *DeviceMonitorServiceImpl) List(c *gin.Context, req *deviceModels.DeviceListReq) (*deviceModels.DeviceInfoListData, error) {
	list, err := d.dao.SelectDeviceList(c, req)
	if err != nil {
		zap.L().Error("device list select error", zap.Error(err))
		return nil, err
	}

	var deviceIds []uint64 = make([]uint64, 0)
	for _, v := range list.Rows {
		deviceIds = append(deviceIds, v.DeviceId)
	}

	var keys []string = make([]string, 0)
	for _, v := range deviceIds {
		keys = append(keys, device.MakeDeviceKey(uint64(v)))
	}

	slice := d.cache.MGet(c, keys)
	for k, v := range slice.Val() {
		str, ok := v.(string)
		if !ok {
			continue
		}
		if str == "" {
			continue
		}
		deviceMetrics := make(map[uint64]map[uint64]*metricModels.DeviceMetricData) // template_id => data_id
		err := json.Unmarshal([]byte(str), &deviceMetrics)
		if err != nil {
			zap.L().Error("json Unmarshal error", zap.Error(err))
			continue
		}
		if deviceMetrics == nil {
			list.Rows[k].Active = false
		} else {
			list.Rows[k].Active = true
		}
		list.Rows[k].TemplateList = deviceMetrics
	}

	// 处理数据
	var dataIds []uint64 = make([]uint64, 0)
	for _, v := range list.Rows {
		for _, templateValue := range v.TemplateList {
			for dataId, _ := range templateValue {
				dataIds = append(dataIds, dataId)
			}
		}
	}

	datas, err := d.deviceConfigDataDao.GetByIds(c, dataIds)
	if err != nil {
		return nil, err
	}

	var dataMap map[uint64]*deviceModels.SysModbusDeviceConfigData = make(map[uint64]*deviceModels.SysModbusDeviceConfigData)
	for _, dataValue := range datas {
		dataMap[uint64(dataValue.DeviceConfigID)] = dataValue
	}

	fmt.Println(dataMap)
	for k, v := range list.Rows {
		for templateId, templateValue := range v.TemplateList {
			for dataId, _ := range templateValue {
				if list.Rows[k].TemplateList == nil {
					continue
				}
				_, ok := list.Rows[k].TemplateList[templateId]
				if !ok {
					continue
				}
				_, ok = list.Rows[k].TemplateList[templateId][dataId]
				if !ok {
					continue
				}
				dataValue, ok := dataMap[dataId]
				if !ok {
					continue
				}
				list.Rows[k].TemplateList[templateId][dataId].Name = dataValue.Name
				list.Rows[k].TemplateList[templateId][dataId].Unit = dataValue.Unit
			}
		}
	}
	return list, nil
}

func (d *DeviceMonitorServiceImpl) Metric(c *gin.Context, req *metricModels.MetricQueryReq) (*metricModels.MetricQueryData, error) {
	info, err := d.deviceConfigDataDao.GetById(c, req.DataId)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, errors.New("数据不存在")
	}
	data, err := d.metricDao.Metric(c, req)
	if err != nil {
		return data, err
	}

	if data == nil {
		return &metricModels.MetricQueryData{
			Labels: make(map[string]string),
			Values: make([]metricModels.MetricQueryValue, 0),
		}, nil
	}

	data.Labels["unit"] = info.Unit
	data.Labels["name"] = info.Name
	return data, nil
}

func (d *DeviceMonitorServiceImpl) Predict(c *gin.Context, req *metricModels.MetricQueryReq) (*metricModels.MetricQueryData, error) {
	info, err := d.deviceConfigDataDao.GetById(c, req.DataId)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, errors.New("数据不存在")
	}
	data, err := d.metricDao.Predict(c, int64(req.DeviceId), info, req)
	if err != nil {
		return data, err
	}

	if data == nil {
		return &metricModels.MetricQueryData{
			Labels: make(map[string]string),
			Values: make([]metricModels.MetricQueryValue, 0),
		}, nil
	}

	data.Labels["unit"] = info.Unit
	data.Labels["name"] = info.Name
	return data, nil
}
