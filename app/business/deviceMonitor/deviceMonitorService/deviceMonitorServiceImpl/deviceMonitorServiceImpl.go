package deviceMonitorServiceImpl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/building/buildingDao"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService"
	"nova-factory-server/app/business/metric/device/metricDao"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/constant/device"
	"nova-factory-server/app/constant/iotdb"
	"nova-factory-server/app/datasource/cache"
	"strconv"
)

type DeviceMonitorServiceImpl struct {
	dao                 deviceDao.IDeviceDao
	metricDao           metricDao.IMetricDao
	cache               cache.Cache
	deviceConfigDataDao deviceDao.ISysModbusDeviceConfigDataDao
	floorDao            buildingDao.FloorDao
	deviceService       deviceService.IDeviceService
}

func NewDeviceMonitorServiceImpl(dao deviceDao.IDeviceDao, cache cache.Cache,
	metricDao metricDao.IMetricDao,
	deviceConfigDataDao deviceDao.ISysModbusDeviceConfigDataDao,
	floorDao buildingDao.FloorDao, deviceService deviceService.IDeviceService) deviceMonitorService.DeviceMonitorService {
	return &DeviceMonitorServiceImpl{
		dao:                 dao,
		cache:               cache,
		metricDao:           metricDao,
		deviceConfigDataDao: deviceConfigDataDao,
		floorDao:            floorDao,
		deviceService:       deviceService,
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
				list.Rows[k].TemplateList[templateId][dataId].GraphEnable = *dataValue.GraphEnable
				list.Rows[k].TemplateList[templateId][dataId].PredictEnable = *dataValue.PredictEnable
				list.Rows[k].TemplateList[templateId][dataId].DataId = uint64(dataValue.DeviceConfigID)
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

func (d *DeviceMonitorServiceImpl) PredictQuery(c *gin.Context, req *metricModels.MetricDataQueryReq) (*metricModels.MetricQueryData, error) {
	if len(req.QueryMetric) != 0 {
		name := ""
		index := 0
		for _, metric := range req.QueryMetric {
			index++
			str := iotdb.MakeDeviceTemplateName(metric.DeviceId, metric.TemplateId, metric.DataId)
			name += str
			if index != len(req.QueryMetric) {
				name += ","
			}
		}
		req.Name = name
	}
	data, err := d.metricDao.Query(c, req)
	if err != nil {
		return data, err
	}

	if data == nil {
		return &metricModels.MetricQueryData{
			Labels: make(map[string]string),
			Values: make([]metricModels.MetricQueryValue, 0),
		}, nil
	}

	return data, nil
}

func (d *DeviceMonitorServiceImpl) DeviceLayout(c *gin.Context, floorId int64) (*deviceMonitorModel.DeviceLayoutData, error) {
	info, err := d.floorDao.Info(c, floorId)
	if err != nil {
		zap.L().Error("get floor info error", zap.Error(err))
		return nil, err
	}

	if info == nil {
		return &deviceMonitorModel.DeviceLayoutData{}, nil
	}

	// 读取设备列表
	if info.LayoutData == nil {
		return &deviceMonitorModel.DeviceLayoutData{}, nil
	}

	var deviceIdMap map[int64]bool = make(map[int64]bool)
	for _, zone := range info.LayoutData.Zones {
		for _, dev := range zone.Devices {
			deviceIdMap[dev.DeviceId] = true
		}
	}

	var deviceIds []int64 = make([]int64, 0)
	for deviceId, _ := range deviceIdMap {
		deviceIds = append(deviceIds, deviceId)
	}

	// 读取设备列表
	list, err := d.deviceService.GetDeviceInfoByIds(c, deviceIds)
	if err != nil {
		return nil, err
	}

	data := &deviceMonitorModel.DeviceLayoutData{}
	fmt.Println(len(list))
	data.Layout = info
	data.DeviceMap = make(map[string]*deviceModels.DeviceVO)
	for _, deviceInfo := range list {
		idStr := strconv.FormatUint(deviceInfo.DeviceId, 10)
		data.DeviceMap[idStr] = deviceInfo
	}
	return data, nil
}
