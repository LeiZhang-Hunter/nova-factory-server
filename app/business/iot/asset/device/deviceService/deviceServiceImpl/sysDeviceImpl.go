package deviceServiceImpl

import (
	"encoding/json"
	"nova-factory-server/app/baize"
	deviceDao2 "nova-factory-server/app/business/iot/asset/device/deviceDao"
	deviceModels2 "nova-factory-server/app/business/iot/asset/device/deviceModels"
	"nova-factory-server/app/business/iot/asset/device/deviceService"
	"nova-factory-server/app/business/iot/deviceMonitor/deviceMonitorDao"
	"nova-factory-server/app/business/iot/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/iot/metric/device/metricDao"
	"nova-factory-server/app/business/iot/metric/device/metricModels"
	"nova-factory-server/app/business/system/systemDao"
	"nova-factory-server/app/business/system/systemModels"
	"nova-factory-server/app/constant/device"
	"nova-factory-server/app/constant/iotdb"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/math"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DeviceService struct {
	iDeviceDao          deviceDao2.IDeviceDao
	iDeviceGroupDao     deviceDao2.IDeviceGroupDao
	iUserDao            systemDao.IUserDao
	metricDao           metricDao.IMetricDao
	dataDao             deviceDao2.ISysModbusDeviceConfigDataDao
	mapDao              deviceMonitorDao.IDeviceDataReportDao
	cache               cache.Cache
	deviceConfigDataDao deviceDao2.ISysModbusDeviceConfigDataDao
}

func NewDeviceService(iDeviceDao deviceDao2.IDeviceDao, iDeviceGroupDao deviceDao2.IDeviceGroupDao,
	iUserDao systemDao.IUserDao, metricDao metricDao.IMetricDao, dataDao deviceDao2.ISysModbusDeviceConfigDataDao,
	mapDao deviceMonitorDao.IDeviceDataReportDao, cache cache.Cache) deviceService.IDeviceService {
	return &DeviceService{
		iDeviceDao:      iDeviceDao,
		iDeviceGroupDao: iDeviceGroupDao,
		iUserDao:        iUserDao,
		metricDao:       metricDao,
		dataDao:         dataDao,
		mapDao:          mapDao,
		cache:           cache,
	}
}

func (d *DeviceService) InsertDevice(c *gin.Context, job *deviceModels2.DeviceInfo) (*deviceModels2.DeviceVO, error) {
	vo, err := d.iDeviceDao.InsertDevice(c, job)
	if err != nil {
		return nil, err
	}

	d.installTable(c, job)
	return vo, nil
}

func (d *DeviceService) installTable(c *gin.Context, info *deviceModels2.DeviceInfo) {
	list, err := d.dataDao.GetByTemplateIds(c, []uint64{info.DeviceProtocolId})
	if err != nil {
		return
	}

	for _, v := range list {
		devKey := iotdb.MakeDeviceTemplateName(int64(info.DeviceId), int64(info.DeviceProtocolId), v.DeviceConfigID)
		err := d.mapDao.Save(c, &deviceMonitorModel.SysIotDbDevMap{
			DeviceID:   int64(info.DeviceId),
			TemplateID: int64(info.DeviceProtocolId),
			DataID:     v.DeviceConfigID,
			Device:     devKey,
			DataName:   v.Name,
			Unit:       v.Unit,
		})
		if err != nil {
			zap.L().Error("save iotdb device map error", zap.Error(err))
		}
		err = d.metricDao.InstallDevice(c, int64(info.DeviceId), v)
		if err != nil {
			zap.L().Error("InstallDevice error", zap.Error(err))
		}
	}
}

func (d *DeviceService) UpdateDevice(c *gin.Context, job *deviceModels2.DeviceInfo) (*deviceModels2.DeviceVO, error) {
	device, err := d.iDeviceDao.UpdateDevice(c, job)
	if err != nil {
		return nil, err
	}
	d.installTable(c, job)
	return device, nil
}

func (d *DeviceService) SelectDeviceList(c *gin.Context, req *deviceModels2.DeviceListReq) (*deviceModels2.DeviceInfoListValue, error) {
	list, err := d.iDeviceDao.SelectDeviceList(c, req)
	if err != nil {
		zap.L().Error("读取列表衰退", zap.Error(err))
		return &deviceModels2.DeviceInfoListValue{
			Rows:  make([]*deviceModels2.DeviceValue, 0),
			Total: 0,
		}, err
	}

	if len(list.Rows) == 0 {
		return &deviceModels2.DeviceInfoListValue{
			Rows:  make([]*deviceModels2.DeviceValue, 0),
			Total: 0,
		}, nil
	}

	// 读取分组id集合
	groupIdMap := make(map[uint64]bool)
	for _, v := range list.Rows {
		if v.DeviceGroupId > 0 {
			groupIdMap[v.DeviceGroupId] = true
		}
	}
	// 格式化服务id
	groupIds := make([]uint64, 0)
	for k, _ := range groupIdMap {
		if k > 0 {
			groupIds = append(groupIds, k)
		}
	}

	//  读取用户id集合
	userIdMap := make(map[int64]bool)
	for _, v := range list.Rows {
		if v.CreateBy > 0 {
			userIdMap[v.CreateBy] = true
		}

		if v.UpdateBy > 0 {
			userIdMap[v.UpdateBy] = true
		}
	}

	// 格式化服务id
	userIds := make([]int64, 0)
	for k, _ := range userIdMap {
		if k > 0 {
			userIds = append(userIds, k)
		}
	}

	// 读取分组列表
	groupList, _ := d.iDeviceGroupDao.GetDeviceGroupByIds(c, groupIds)
	groupVoMap := make(map[uint64]*deviceModels2.DeviceGroupVO)
	for _, v := range groupList {
		groupVoMap[v.GroupId] = v
	}

	users := d.iUserDao.SelectByUserIds(c, userIds)
	userVoMap := make(map[int64]*systemModels.SysUserDML)
	for _, v := range users {
		userVoMap[v.UserId] = v
	}

	ret := make([]*deviceModels2.DeviceValue, 0)
	for _, v := range list.Rows {
		var actions []string = make([]string, 0)
		if len(v.Action) != 0 {
			json.Unmarshal([]byte((v.Action)), &actions)
		}

		var groupName string
		groupVo, ok := groupVoMap[v.DeviceGroupId]
		if ok {
			groupName = groupVo.Name
		}

		var createUserName string
		var updateUserName string
		userVo, ok := userVoMap[v.CreateBy]
		if ok {
			createUserName = userVo.UserName
		}

		userVo, ok = userVoMap[v.UpdateBy]
		if ok {
			updateUserName = userVo.UserName
		}

		value := &deviceModels2.DeviceValue{
			DeviceId:          v.DeviceId,
			DeviceGroupId:     v.DeviceGroupId,
			DeviceClassId:     v.DeviceClassId,
			DeviceProtocolId:  v.DeviceProtocolId,
			DeviceBuildingId:  v.DeviceBuildingId,
			Name:              *v.Name,
			DeviceGroupName:   groupName,
			CommunicationType: v.CommunicationType,
			ProtocolType:      v.ProtocolType,
			DeviceGatewayID:   v.DeviceGatewayID,
			Number:            *v.Number,
			Type:              *v.Type,
			Action:            actions,
			Extension:         v.Extension,
			ControlType:       v.ControlType,
			CreateUserName:    createUserName,
			Status:            v.Status,
			Enable:            v.Enable,
			ModelURL:          v.ModelURL,
			UpdateUserName:    updateUserName,
			BaseEntity: baize.BaseEntity{
				CreateTime: v.CreateTime,
				UpdateTime: v.UpdateTime,
			},
		}

		ret = append(ret, value)
	}
	return &deviceModels2.DeviceInfoListValue{
		Rows:  ret,
		Total: list.Total,
	}, nil
}

func (d *DeviceService) DeleteByDeviceIds(c *gin.Context, ids []int64) error {
	err := d.iDeviceDao.DeleteByDeviceIds(c, ids)
	if err != nil {
		return err
	}
	data, err := d.iDeviceDao.GetByIds(c, ids)
	if err != nil {
		zap.L().Error("get by ids error", zap.Error(err))
		return err
	}
	for _, value := range data {
		d.unInstallTable(c, value)
	}
	return nil
}

func (d *DeviceService) unInstallTable(c *gin.Context, info *deviceModels2.DeviceVO) {
	list, err := d.dataDao.GetByTemplateIds(c, []uint64{info.DeviceProtocolId})
	if err != nil {
		return
	}

	for _, v := range list {
		devKey := iotdb.MakeDeviceTemplateName(int64(info.DeviceId), int64(info.DeviceProtocolId), v.DeviceConfigID)
		err = d.mapDao.Remove(c, devKey)
		if err != nil {
			zap.L().Error("uninstall device template error", zap.Error(err))
		}
		err := d.metricDao.UnInStallDevice(c, int64(info.DeviceId), v.TemplateID, v.DeviceConfigID)
		if err != nil {
			zap.L().Error("InstallDevice error", zap.Error(err))
		}
	}
}

func (d *DeviceService) GetById(c *gin.Context, id int64) (*deviceModels2.DeviceVO, error) {
	device, err := d.iDeviceDao.GetById(c, id)
	if err != nil {
		zap.L().Error("Get device error", zap.Error(err))
	}
	return device, err
}

func (d *DeviceService) GetMetricByTag(c *gin.Context, req *deviceModels2.DeviceTagListReq) (*deviceModels2.DeviceMetricInfoListValue, error) {
	list, err := d.iDeviceDao.SelectPublicDeviceList(c, &deviceModels2.DeviceListReq{
		Number: &req.Tags,
		BaseEntityDQL: baize.BaseEntityDQL{
			DataScope: req.DataScope,
			OrderBy:   req.OrderBy,
			IsAsc:     req.IsAsc,
			Page:      req.Page,
			Size:      req.Size,
		},
	})
	if err != nil {
		zap.L().Error("读取列表衰退", zap.Error(err))
		return &deviceModels2.DeviceMetricInfoListValue{
			Rows:  make([]*deviceModels2.DeviceVO, 0),
			Total: 0,
		}, err
	}

	if len(list.Rows) == 0 {
		return &deviceModels2.DeviceMetricInfoListValue{
			Rows:  make([]*deviceModels2.DeviceVO, 0),
			Total: 0,
		}, nil
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

	datas, err := d.dataDao.GetByIds(c, dataIds)
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
				list.Rows[k].TemplateList[templateId][dataId].GraphEnable = *dataValue.GraphEnable
				list.Rows[k].TemplateList[templateId][dataId].PredictEnable = *dataValue.PredictEnable
				list.Rows[k].TemplateList[templateId][dataId].DataId = uint64(dataValue.DeviceConfigID)
				list.Rows[k].TemplateList[templateId][dataId].Value = math.RoundFloat(list.Rows[k].TemplateList[templateId][dataId].Value, 2)
				if dataValue.Annotation != "" {
					annotations := make([]deviceModels2.Annotation, 0)
					err := json.Unmarshal([]byte(dataValue.Annotation), &annotations)
					if err != nil {
						zap.L().Error("json Unmarshal error", zap.Error(err))
						continue
					}
					for _, annotationValue := range annotations {
						list.Rows[k].TemplateList[templateId][dataId].Attributes[annotationValue.Key] = annotationValue.Value
					}
				}

			}
		}
	}
	return &deviceModels2.DeviceMetricInfoListValue{
		Rows:  list.Rows,
		Total: list.Total,
	}, nil
}

func (d *DeviceService) GetDeviceInfoByIds(c *gin.Context, ids []int64) ([]*deviceModels2.DeviceVO, error) {
	if len(ids) == 0 {
		return make([]*deviceModels2.DeviceVO, 0), nil
	}
	deviceList, err := d.iDeviceDao.GetByIds(c, ids)
	if err != nil {
		zap.L().Error("Get device error", zap.Error(err))
	}

	var keys []string = make([]string, 0)
	for _, v := range deviceList {
		keys = append(keys, device.MakeDeviceKey(uint64(v.DeviceId)))
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
			deviceList[k].Active = false
		} else {
			deviceList[k].Active = true
		}
	}

	return deviceList, err
}

// StatCount 统计设备，在线、不在线、异常
func (d *DeviceService) StatCount(c *gin.Context) (*deviceModels2.DeviceStatData, error) {
	var data deviceModels2.DeviceStatData
	list, err := d.iDeviceDao.All(c)
	if err != nil {
		return &data, nil
	}

	data.Total = int64(len(list))
	var keys []string = make([]string, 0)
	for _, v := range list {
		keys = append(keys, device.MakeDeviceKey(v.DeviceId))
	}

	slice := d.cache.MGet(c, keys)
	for k, v := range slice.Val() {

		if list[k].Status == device.WORK_EXCEPTION {
			data.Exception++
			continue
		}

		if list[k].Status == device.WORK_MAINTENANCE {
			data.Maintenance++
			continue
		}

		str, ok := v.(string)
		if !ok {
			data.OffLine++
			continue
		}
		if str == "" {
			data.OffLine++
			continue
		}
		deviceMetrics := make(map[uint64]map[uint64]*metricModels.DeviceMetricData) // template_id => data_id
		err := json.Unmarshal([]byte(str), &deviceMetrics)
		if err != nil {
			data.Exception++
			zap.L().Error("json Unmarshal error", zap.Error(err))
			continue
		}
		if deviceMetrics == nil {
			data.OffLine++
		} else {
			data.Online++
		}
	}
	return &data, nil
}
