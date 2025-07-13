package deviceServiceImpl

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorDao"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/metric/device/metricDao"
	"nova-factory-server/app/business/system/systemDao"
	"nova-factory-server/app/business/system/systemModels"
	"nova-factory-server/app/constant/iotdb"
)

type DeviceService struct {
	iDeviceDao      deviceDao.IDeviceDao
	iDeviceGroupDao deviceDao.IDeviceGroupDao
	iUserDao        systemDao.IUserDao
	metricDao       metricDao.IMetricDao
	dataDao         deviceDao.ISysModbusDeviceConfigDataDao
	mapDao          deviceMonitorDao.IDeviceDataReportDao
}

func NewDeviceService(iDeviceDao deviceDao.IDeviceDao, iDeviceGroupDao deviceDao.IDeviceGroupDao,
	iUserDao systemDao.IUserDao, metricDao metricDao.IMetricDao, dataDao deviceDao.ISysModbusDeviceConfigDataDao,
	mapDao deviceMonitorDao.IDeviceDataReportDao) deviceService.IDeviceService {
	return &DeviceService{
		iDeviceDao:      iDeviceDao,
		iDeviceGroupDao: iDeviceGroupDao,
		iUserDao:        iUserDao,
		metricDao:       metricDao,
		dataDao:         dataDao,
		mapDao:          mapDao,
	}
}

func (d *DeviceService) InsertDevice(c *gin.Context, job *deviceModels.DeviceInfo) (*deviceModels.DeviceVO, error) {
	vo, err := d.iDeviceDao.InsertDevice(c, job)
	if err != nil {
		return nil, err
	}

	d.installTable(c, job)
	return vo, nil
}

func (d *DeviceService) installTable(c *gin.Context, info *deviceModels.DeviceInfo) {
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

func (d *DeviceService) UpdateDevice(c *gin.Context, job *deviceModels.DeviceInfo) (*deviceModels.DeviceVO, error) {
	device, err := d.iDeviceDao.UpdateDevice(c, job)
	if err != nil {
		return nil, err
	}
	d.installTable(c, job)
	return device, nil
}

func (d *DeviceService) SelectDeviceList(c *gin.Context, req *deviceModels.DeviceListReq) (*deviceModels.DeviceInfoListValue, error) {
	list, err := d.iDeviceDao.SelectDeviceList(c, req)
	if err != nil {
		zap.L().Error("读取列表衰退", zap.Error(err))
		return &deviceModels.DeviceInfoListValue{
			Rows:  make([]*deviceModels.DeviceValue, 0),
			Total: 0,
		}, err
	}

	if len(list.Rows) == 0 {
		return &deviceModels.DeviceInfoListValue{
			Rows:  make([]*deviceModels.DeviceValue, 0),
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
	groupVoMap := make(map[uint64]*deviceModels.DeviceGroupVO)
	for _, v := range groupList {
		groupVoMap[v.GroupId] = v
	}

	users := d.iUserDao.SelectByUserIds(c, userIds)
	userVoMap := make(map[int64]*systemModels.SysUserDML)
	for _, v := range users {
		userVoMap[v.UserId] = v
	}

	ret := make([]*deviceModels.DeviceValue, 0)
	for _, v := range list.Rows {
		var actions []string = make([]string, 0)
		if len(v.Action) != 0 {
			json.Unmarshal([]byte((v.Action)), &actions)
		}

		var extensions map[string]string = make(map[string]string)
		if len(v.Extension) != 0 {
			json.Unmarshal([]byte((v.Extension)), &extensions)
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

		value := &deviceModels.DeviceValue{
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
			UpdateUserName:    updateUserName,
			BaseEntity: baize.BaseEntity{
				CreateTime: v.CreateTime,
				UpdateTime: v.UpdateTime,
			},
		}

		ret = append(ret, value)
	}
	return &deviceModels.DeviceInfoListValue{
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

func (d *DeviceService) unInstallTable(c *gin.Context, info *deviceModels.DeviceVO) {
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
