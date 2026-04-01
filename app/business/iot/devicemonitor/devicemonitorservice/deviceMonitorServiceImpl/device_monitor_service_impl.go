package deviceMonitorServiceImpl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	buildingDao2 "nova-factory-server/app/business/iot/asset/building/buildingdao"
	deviceDao2 "nova-factory-server/app/business/iot/asset/device/devicedao"
	deviceModels2 "nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/asset/device/deviceservice"
	deviceMonitorModel2 "nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	deviceMonitorService2 "nova-factory-server/app/business/iot/devicemonitor/devicemonitorservice"
	metricDao2 "nova-factory-server/app/business/iot/metric/device/metricdao"
	metricModels2 "nova-factory-server/app/business/iot/metric/device/metricmodels"
	"nova-factory-server/app/constant/device"
	"nova-factory-server/app/constant/iotdb"
	"nova-factory-server/app/datasource/cache"
	"strconv"
	"time"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/os/gtime"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	controlService "github.com/novawatcher-io/nova-factory-payload/control/v1"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DeviceMonitorServiceImpl struct {
	dao                 deviceDao2.IDeviceDao
	metricDao           metricDao2.IMetricDao
	cache               cache.Cache
	deviceConfigDataDao deviceDao2.ISysModbusDeviceConfigDataDao
	floorDao            buildingDao2.FloorDao
	deviceService       deviceservice.IDeviceService
	deviceControl       deviceMonitorService2.DeviceControlService
	controlLogDao       metricDao2.IControlLogDao
	buildingDao         buildingDao2.BuildingDao
}

func NewDeviceMonitorServiceImpl(dao deviceDao2.IDeviceDao, cache cache.Cache,
	metricDao metricDao2.IMetricDao,
	deviceConfigDataDao deviceDao2.ISysModbusDeviceConfigDataDao,
	floorDao buildingDao2.FloorDao, deviceService deviceservice.IDeviceService,
	deviceControl deviceMonitorService2.DeviceControlService,
	controlLogDao metricDao2.IControlLogDao,
	buildingDao buildingDao2.BuildingDao) deviceMonitorService2.DeviceMonitorService {
	return &DeviceMonitorServiceImpl{
		dao:                 dao,
		cache:               cache,
		metricDao:           metricDao,
		deviceConfigDataDao: deviceConfigDataDao,
		floorDao:            floorDao,
		deviceService:       deviceService,
		deviceControl:       deviceControl,
		controlLogDao:       controlLogDao,
		buildingDao:         buildingDao,
	}
}

func (d *DeviceMonitorServiceImpl) List(c *gin.Context, req *deviceModels2.DeviceListReq) (*deviceModels2.DeviceInfoListData, error) {
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
		deviceMetrics := make(map[uint64]map[uint64]*metricModels2.DeviceMetricData) // template_id => data_id
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

	var dataMap map[uint64]*deviceModels2.SysModbusDeviceConfigData = make(map[uint64]*deviceModels2.SysModbusDeviceConfigData)
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
				list.Rows[k].TemplateList[templateId][dataId].DataFormat = dataValue.DataFormat
				list.Rows[k].TemplateList[templateId][dataId].GraphEnable = *dataValue.GraphEnable
				list.Rows[k].TemplateList[templateId][dataId].PredictEnable = *dataValue.PredictEnable
				list.Rows[k].TemplateList[templateId][dataId].Mode = dataValue.Mode
				list.Rows[k].TemplateList[templateId][dataId].DataType = dataValue.DataType
				list.Rows[k].TemplateList[templateId][dataId].Type = dataValue.Type
				list.Rows[k].TemplateList[templateId][dataId].DataId = uint64(dataValue.DeviceConfigID)
				list.Rows[k].TemplateList[templateId][dataId].ConfigurationEnable = dataValue.ConfigurationEnable
			}
		}
	}
	// 补全建筑物名称
	buildingIds := make([]uint64, 0)
	buildingIdMap := make(map[uint64]bool)
	for _, row := range list.Rows {
		if row.DeviceBuildingId > 0 && !buildingIdMap[row.DeviceBuildingId] {
			buildingIds = append(buildingIds, row.DeviceBuildingId)
			buildingIdMap[row.DeviceBuildingId] = true
		}
	}

	if len(buildingIds) > 0 {
		buildings, err := d.buildingDao.GetByIds(c, buildingIds)
		if err != nil {
			zap.L().Error("get buildings error", zap.Error(err))
		} else {
			buildingNameMap := make(map[uint64]string)
			for _, b := range buildings {
				buildingNameMap[uint64(b.ID)] = b.Name
			}
			for i, row := range list.Rows {
				if name, ok := buildingNameMap[row.DeviceBuildingId]; ok {
					list.Rows[i].BuildingName = name
				}
			}
		}
	}

	return list, nil
}

func (d *DeviceMonitorServiceImpl) Metric(c *gin.Context, req *metricModels2.MetricQueryReq) (*metricModels2.MetricQueryData, error) {
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
		return &metricModels2.MetricQueryData{
			Labels: make(map[string]string),
			Values: make([]metricModels2.MetricQueryValue, 0),
		}, nil
	}

	data.Labels["unit"] = info.Unit
	data.Labels["name"] = info.Name
	return data, nil
}

func (d *DeviceMonitorServiceImpl) Predict(c *gin.Context, req *metricModels2.MetricQueryReq) (*metricModels2.MetricQueryData, error) {
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
		return &metricModels2.MetricQueryData{
			Labels: make(map[string]string),
			Values: make([]metricModels2.MetricQueryValue, 0),
		}, nil
	}

	data.Labels["unit"] = info.Unit
	data.Labels["name"] = info.Name
	return data, nil
}

func (d *DeviceMonitorServiceImpl) PredictQuery(c *gin.Context, req *metricModels2.MetricDataQueryReq) (*metricModels2.MetricQueryData, error) {
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
		return &metricModels2.MetricQueryData{
			Labels: make(map[string]string),
			Values: make([]metricModels2.MetricQueryValue, 0),
		}, nil
	}

	return data, nil
}

func (d *DeviceMonitorServiceImpl) DeviceLayout(c *gin.Context, floorId int64) (*deviceMonitorModel2.DeviceLayoutData, error) {
	info, err := d.floorDao.Info(c, floorId)
	if err != nil {
		zap.L().Error("get floor info error", zap.Error(err))
		return nil, err
	}

	if info == nil {
		return &deviceMonitorModel2.DeviceLayoutData{}, nil
	}

	// 读取设备列表
	if info.LayoutData == nil {
		return &deviceMonitorModel2.DeviceLayoutData{}, nil
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

	data := &deviceMonitorModel2.DeviceLayoutData{}
	data.Layout = info
	data.DeviceMap = make(map[string]*deviceModels2.DeviceVO)
	for _, deviceInfo := range list {
		idStr := strconv.FormatUint(deviceInfo.DeviceId, 10)
		deviceInfo.ExtensionInfo = nil
		deviceInfo.Extension = ""
		data.DeviceMap[idStr] = deviceInfo
	}
	return data, nil
}

func (d *DeviceMonitorServiceImpl) ControlStatus(c context.Context, req *deviceMonitorModel2.ControlStatusReq) (*deviceMonitorModel2.ControlStatusRes, error) {
	if len(req.Items) == 0 {
		return &deviceMonitorModel2.ControlStatusRes{Items: []deviceMonitorModel2.ControlStatusItemRes{}}, nil
	}

	keys := make([]string, 0, len(req.Items))
	for _, item := range req.Items {
		keys = append(keys, fmt.Sprintf(device.DEVICE_CONTROL_KEY, item.DeviceId, item.DataId))
	}

	sliceCmd := d.cache.MGet(c, keys)
	vals, err := sliceCmd.Result()
	if err != nil {
		zap.L().Error("MGet error", zap.Error(err))
		return nil, err
	}

	resItems := make([]deviceMonitorModel2.ControlStatusItemRes, 0, len(req.Items))
	for i, val := range vals {
		status := 0
		if val != nil {
			status = 1 // in progress
		}
		resItems = append(resItems, deviceMonitorModel2.ControlStatusItemRes{
			DeviceId: req.Items[i].DeviceId,
			DataId:   req.Items[i].DataId,
			Status:   status,
		})
	}

	return &deviceMonitorModel2.ControlStatusRes{
		Items: resItems,
	}, nil
}

func (d *DeviceMonitorServiceImpl) Control(c *gin.Context, req *deviceMonitorModel2.ControlReq) (*deviceMonitorModel2.ControlRes, error) {
	requestId := uuid.New().String()
	value, err := req.Value.ToValue()
	if err != nil {
		zap.L().Error("ToValue error", zap.Error(err))
		return nil, err
	}

	// 读取设备信息
	deviceInfo, err := d.dao.GetById(c, int64(req.DeviceId))
	if err != nil {
		zap.L().Error("get device error", zap.Error(err))
		return nil, errors.New("读取设备信息错误")
	}
	if deviceInfo == nil {
		return nil, errors.New("设备不存在")
	}

	var deviceName string
	if deviceInfo.Name != nil {
		deviceName = *deviceInfo.Name
	}

	dataInfo, err := d.deviceConfigDataDao.GetById(c, req.DataId)
	if err != nil {
		zap.L().Error("get device data error", zap.Error(err))
		return nil, errors.New("读取设备模板错误")
	}

	if dataInfo == nil {
		zap.L().Error("get device data error", zap.Error(err))
		return nil, errors.New("设备模板信息不存在")
	}

	content, err := json.Marshal(req)
	if err != nil {
		return nil, errors.New("下发内容编码失败")
	}

	serieId := time.Now().UnixNano()
	var now *gtime.Time = gtime.Now()
	err = d.controlLogDao.Export(context.Background(), []*metricModels2.NovaControlLog{
		{
			DeviceId:      req.DeviceId,
			DeviceName:    deviceName,
			DataId:        req.DataId,
			DataName:      dataInfo.Name,
			Message:       fmt.Sprintf("[下发中]:控制<%s><%s>指令 value %s 下发中", deviceName, dataInfo.Name, string(content)),
			Type:          "manual",
			SeriesId:      uint64(serieId),
			Attributes:    make(map[string]string),
			StartTimeUnix: now,
			TimeUnix:      now,
		},
	})
	if err != nil {
		zap.L().Error("Export error", zap.Error(err))
		return nil, err
	}
	res, err := d.broadcastControlToServers(c, d.newControlBroadcastRequest(requestId, req, value))

	if err != nil {
		zap.L().Error("BroadcastControl error", zap.Error(err))
		return nil, err
	}

	if res.Code == 201 {
		return &deviceMonitorModel2.ControlRes{
			Code: 0,
			Msg:  "success",
		}, nil
	}

	if res.Code != 0 {
		return &deviceMonitorModel2.ControlRes{
			Code: int(res.Code),
			Msg:  res.Message,
		}, nil
	}

	return &deviceMonitorModel2.ControlRes{
		Code: 0,
		Msg:  "success",
	}, nil
}

func (d *DeviceMonitorServiceImpl) newControlBroadcastRequest(requestId string, req *deviceMonitorModel2.ControlReq, value *controlService.Value) *controlService.ControlRequest {
	return &controlService.ControlRequest{
		RequestId: requestId,
		DeviceId:  req.DeviceId,
		AgentId:   req.AgentId,
		DataId:    req.DataId,
		Value:     value,
		Timestamp: timestamppb.Now(),
	}
}

func (d *DeviceMonitorServiceImpl) broadcastControlToServers(ctx context.Context, req *controlService.ControlRequest) (*controlService.ControlResponse, error) {
	addressList := viper.GetStringSlice("daemonize.server_list")
	if len(addressList) == 0 {
		return &controlService.ControlResponse{
			Code:    -1,
			Message: "地址没有配置",
		}, nil
	}
	for _, address := range addressList {
		client := controlService.NewControlServiceClient(grpcx.Client.MustNewGrpcClientConn(address))
		remoteRes, remoteErr := client.BroadcastControl(ctx, req)
		if remoteErr != nil {
			zap.L().Warn("remote broadcast control error", zap.String("address", address), zap.Error(remoteErr))
			continue
		}
		return remoteRes, nil
	}
	return &controlService.ControlResponse{
		Code:    0,
		Message: "成功",
	}, nil
}
func (d *DeviceMonitorServiceImpl) GetRealTimeInfo(c *gin.Context, deviceId uint64) (*deviceModels2.DeviceVO, error) {
	// 1. 读取设备信息
	deviceInfo, err := d.dao.GetById(c, int64(deviceId))
	if err != nil {
		zap.L().Error("get device error", zap.Error(err))
		return nil, errors.New("读取设备信息错误")
	}
	if deviceInfo == nil {
		return nil, errors.New("设备不存在")
	}

	vo := deviceInfo

	// 2. 从缓存读取实时指标
	key := device.MakeDeviceKey(deviceId)
	val, err := d.cache.Get(c, key)
	if err != nil {
		// 缓存没有数据，不报错，只将 Active 设为 false
		vo.Active = false
		return vo, nil
	}

	if val != "" {
		deviceMetrics := make(map[uint64]map[uint64]*metricModels2.DeviceMetricData)
		err = json.Unmarshal([]byte(val), &deviceMetrics)
		if err != nil {
			zap.L().Error("json Unmarshal error", zap.Error(err))
		} else {
			vo.TemplateList = deviceMetrics
			vo.Active = (deviceMetrics != nil)
		}
	}

	// 3. 补全建筑物名称
	if vo.DeviceBuildingId > 0 {
		buildings, err := d.buildingDao.GetByIds(c, []uint64{vo.DeviceBuildingId})
		if err != nil {
			zap.L().Error("get building error", zap.Error(err))
		} else if len(buildings) > 0 {
			vo.BuildingName = buildings[0].Name
		}
	}

	// 3. 补全指标元数据
	if vo.TemplateList != nil {
		var dataIds []uint64 = make([]uint64, 0)
		for _, templateValue := range vo.TemplateList {
			for dataId := range templateValue {
				dataIds = append(dataIds, dataId)
			}
		}

		if len(dataIds) > 0 {
			datas, err := d.deviceConfigDataDao.GetByIds(c, dataIds)
			if err != nil {
				zap.L().Error("get device config data error", zap.Error(err))
				return vo, nil // 返回基本信息，忽略指标补全错误
			}

			var dataMap = make(map[uint64]*deviceModels2.SysModbusDeviceConfigData)
			for _, dataValue := range datas {
				dataMap[uint64(dataValue.DeviceConfigID)] = dataValue
			}

			for templateId, templateValue := range vo.TemplateList {
				for dataId := range templateValue {
					dataValue, ok := dataMap[dataId]
					if !ok {
						continue
					}
					vo.TemplateList[templateId][dataId].Name = dataValue.Name
					vo.TemplateList[templateId][dataId].Unit = dataValue.Unit
					vo.TemplateList[templateId][dataId].DataFormat = dataValue.DataFormat
					vo.TemplateList[templateId][dataId].GraphEnable = *dataValue.GraphEnable
					vo.TemplateList[templateId][dataId].PredictEnable = *dataValue.PredictEnable
					vo.TemplateList[templateId][dataId].Mode = dataValue.Mode
					vo.TemplateList[templateId][dataId].DataType = dataValue.DataType
					vo.TemplateList[templateId][dataId].Type = dataValue.Type
					vo.TemplateList[templateId][dataId].DataId = uint64(dataValue.DeviceConfigID)
					vo.TemplateList[templateId][dataId].ConfigurationEnable = dataValue.ConfigurationEnable
				}
			}
		}
	}

	return vo, nil
}

func (d *DeviceMonitorServiceImpl) GetRealTimeInfoList(c *gin.Context, deviceIds []string) ([]*deviceModels2.DeviceVO, error) {
	if len(deviceIds) == 0 {
		return make([]*deviceModels2.DeviceVO, 0), nil
	}

	var ids []int64
	for _, id := range deviceIds {
		v, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			zap.L().Error("strconv.ParseInt error", zap.Error(err))
		}
		ids = append(ids, int64(v))
	}

	// 1. 读取设备信息
	deviceInfos, err := d.dao.GetByIds(c, ids)
	if err != nil {
		zap.L().Error("get devices error", zap.Error(err))
		return nil, errors.New("批量读取设备信息错误")
	}

	var voList []*deviceModels2.DeviceVO
	var keys []string
	var buildingIds []uint64
	buildingIdMap := make(map[uint64]bool)

	for _, deviceInfo := range deviceInfos {
		voList = append(voList, deviceInfo)
		keys = append(keys, device.MakeDeviceKey(deviceInfo.DeviceId))

		if deviceInfo.DeviceBuildingId > 0 && !buildingIdMap[deviceInfo.DeviceBuildingId] {
			buildingIds = append(buildingIds, deviceInfo.DeviceBuildingId)
			buildingIdMap[deviceInfo.DeviceBuildingId] = true
		}
	}

	// 2. 从缓存读取实时指标
	sliceCmd := d.cache.MGet(c, keys)
	vals, err := sliceCmd.Result()
	if err != nil {
		zap.L().Error("mget cache error", zap.Error(err))
		for _, vo := range voList {
			vo.Active = false
		}
	} else {
		for i, val := range vals {
			vo := voList[i]
			str, ok := val.(string)
			if ok && str != "" {
				deviceMetrics := make(map[uint64]map[uint64]*metricModels2.DeviceMetricData)
				err = json.Unmarshal([]byte(str), &deviceMetrics)
				if err != nil {
					zap.L().Error("json Unmarshal error", zap.Error(err))
					vo.Active = false
				} else {
					vo.TemplateList = deviceMetrics
					vo.Active = (deviceMetrics != nil)
				}
			} else {
				vo.Active = false
			}
		}
	}

	// 3. 补全建筑物名称
	if len(buildingIds) > 0 {
		buildings, err := d.buildingDao.GetByIds(c, buildingIds)
		if err != nil {
			zap.L().Error("get buildings error", zap.Error(err))
		} else {
			buildingNameMap := make(map[uint64]string)
			for _, b := range buildings {
				buildingNameMap[uint64(b.ID)] = b.Name
			}
			for _, vo := range voList {
				if name, ok := buildingNameMap[vo.DeviceBuildingId]; ok {
					vo.BuildingName = name
				}
			}
		}
	}

	// 4. 补全指标元数据
	var allDataIds []uint64
	dataIdMap := make(map[uint64]bool)
	for _, vo := range voList {
		if vo.TemplateList != nil {
			for _, templateValue := range vo.TemplateList {
				for dataId := range templateValue {
					if !dataIdMap[dataId] {
						allDataIds = append(allDataIds, dataId)
						dataIdMap[dataId] = true
					}
				}
			}
		}
	}

	if len(allDataIds) > 0 {
		datas, err := d.deviceConfigDataDao.GetByIds(c, allDataIds)
		if err != nil {
			zap.L().Error("get device config data error", zap.Error(err))
			return voList, nil // 返回基本信息，忽略指标补全错误
		}

		var dataMap = make(map[uint64]*deviceModels2.SysModbusDeviceConfigData)
		for _, dataValue := range datas {
			dataMap[uint64(dataValue.DeviceConfigID)] = dataValue
		}

		for _, vo := range voList {
			if vo.TemplateList != nil {
				for templateId, templateValue := range vo.TemplateList {
					for dataId := range templateValue {
						dataValue, ok := dataMap[dataId]
						if !ok {
							continue
						}
						vo.TemplateList[templateId][dataId].Name = dataValue.Name
						vo.TemplateList[templateId][dataId].Unit = dataValue.Unit
						vo.TemplateList[templateId][dataId].DataFormat = dataValue.DataFormat
						vo.TemplateList[templateId][dataId].GraphEnable = *dataValue.GraphEnable
						vo.TemplateList[templateId][dataId].PredictEnable = *dataValue.PredictEnable
						vo.TemplateList[templateId][dataId].Mode = dataValue.Mode
						vo.TemplateList[templateId][dataId].DataType = dataValue.DataType
						vo.TemplateList[templateId][dataId].Type = dataValue.Type
						vo.TemplateList[templateId][dataId].DataId = uint64(dataValue.DeviceConfigID)
						vo.TemplateList[templateId][dataId].ConfigurationEnable = dataValue.ConfigurationEnable
					}
				}
			}
		}
	}

	return voList, nil
}
