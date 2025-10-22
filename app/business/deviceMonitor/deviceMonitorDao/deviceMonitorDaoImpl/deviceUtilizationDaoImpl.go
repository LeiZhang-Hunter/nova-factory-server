package deviceMonitorDaoImpl

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/building/buildingDao"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorDao"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/system/systemDao"
	"nova-factory-server/app/constant/device"
	"nova-factory-server/app/datasource/iotdb"
	"nova-factory-server/app/utils/math"
	"nova-factory-server/app/utils/time"
	systime "time"
)

type DeviceUtilizationDaoImpl struct {
	iotDb          *iotdb.IotDb
	deviceDao      deviceDao.IDeviceDao
	deviceBuildDao buildingDao.BuildingDao
	shiftDao       systemDao.ISysShiftDao
}

func NewDeviceUtilizationDaoImpl(iotDb *iotdb.IotDb, shiftDao systemDao.ISysShiftDao, deviceDao deviceDao.IDeviceDao, deviceBuildDao buildingDao.BuildingDao) deviceMonitorDao.DeviceUtilizationDao {
	return &DeviceUtilizationDaoImpl{
		iotDb:          iotDb,
		shiftDao:       shiftDao,
		deviceDao:      deviceDao,
		deviceBuildDao: deviceBuildDao,
	}
}

// statRun 统计运行时设备
func (d *DeviceUtilizationDaoImpl) statDeviceStat(c *gin.Context, startTime string, endTime string, status device.RUN_STATUS) (*deviceMonitorModel.DeviceStatusList, error) {
	session, err := d.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return nil, err
	}
	defer d.iotDb.PutSession(session)

	var timeout int64 = 5000

	sql := fmt.Sprintf("select sum(duration) as value from root.run_status_device.** where time > %s and time < %s and status = %d align by device",
		startTime, endTime, status)

	statement, err := session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		zap.L().Error("ExecuteQueryStatement error", zap.Error(err))
		return nil, err
	}
	data := deviceMonitorModel.NewDeviceStatusList()
	for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
		v := statement.GetDouble("value")
		deviceName := statement.GetText("Device")
		var deviceId int64
		_, err := fmt.Sscanf(deviceName, "root.run_status_device.dev%d", &deviceId)
		if err != nil {
			zap.L().Error("fmt Sscanf error", zap.Error(err))
			continue
		}
		data.List = append(data.List, deviceMonitorModel.DeviceStatus{
			DeviceId: deviceId,
			Value:    v,
			Status:   status,
		})
	}

	return data, nil
}

// Stat 统计稼动率
func (d *DeviceUtilizationDaoImpl) Stat(c *gin.Context, req *deviceMonitorModel.DeviceUtilizationReq) ([]*deviceMonitorModel.DeviceUtilizationData, error) {
	var deviceUtilizationList []*deviceMonitorModel.DeviceUtilizationData = make([]*deviceMonitorModel.DeviceUtilizationData, 0)

	if req == nil {
		return deviceUtilizationList, errors.New("参数错误")
	}

	// 查询所有班次，用来计算工作时间的稼动率
	shift, err := d.shiftDao.GetEnableShift(c)
	if err != nil {
		zap.L().Error("读取班次错误", zap.Error(err))
		return deviceUtilizationList, errors.New("读取班次错误")
	}
	if shift == nil {
		zap.L().Error("读取班次错误")
		return deviceUtilizationList, errors.New("读取班次错误")
	}

	shiftTime := 0
	for _, config := range shift {
		endTime := config.EndTime
		if endTime > 86400 {
			endTime = 86400
		}
		interval := endTime - config.BeginTime
		if interval <= 0 {
			zap.L().Error("班次时间设置错误")
			return deviceUtilizationList, errors.New("班次时间设置错误")
		}

		shiftTime += int(interval)

		// 添加超出当天时间的值
		if config.EndTime > 86400 {
			v := int(config.EndTime) - 86400
			if v < 0 {
				continue
			}
			shiftTime += v
		}
	}

	if shiftTime > 86400 {
		shiftTime = 86400
	}

	if shiftTime == 0 {
		return deviceUtilizationList, nil
	}

	var startTime string
	if req.Start > 0 {
		startTime = time.GetStartTime(req.Start, 200)
	} else {
		start := systime.Now().UnixMilli() - 86400*1000
		startTime = time.GetStartTime(uint64(start), 200)
	}
	endTime := time.GetEndTimeUseNow(req.End, true)

	runList, err := d.statDeviceStat(c, startTime, endTime, device.RUNNING)
	if err != nil {
		return deviceUtilizationList, err
	}
	waitingList, err := d.statDeviceStat(c, startTime, endTime, device.WAITING)
	if err != nil {
		return deviceUtilizationList, err
	}

	// 批量读取设备id
	var deviceRunMap map[int64]*deviceMonitorModel.DeviceStatus = make(map[int64]*deviceMonitorModel.DeviceStatus)
	var deviceWaitMap map[int64]*deviceMonitorModel.DeviceStatus = make(map[int64]*deviceMonitorModel.DeviceStatus)
	var deviceMap map[int64]*deviceModels.DeviceVO = make(map[int64]*deviceModels.DeviceVO)
	for _, v := range runList.List {
		if v.DeviceId == 0 {
			continue
		}
		deviceMap[v.DeviceId] = nil
		deviceRunMap[v.DeviceId] = &v
	}

	for _, v := range waitingList.List {
		if v.DeviceId == 0 {
			continue
		}
		deviceMap[v.DeviceId] = nil
		deviceWaitMap[v.DeviceId] = &v
	}

	// 设备id集合
	var deviceIds []int64 = make([]int64, 0)
	for id, _ := range deviceMap {
		deviceIds = append(deviceIds, id)
	}

	if len(deviceIds) == 0 {
		return deviceUtilizationList, nil
	}

	// 建筑映射表
	var buildingMap map[uint64]uint64 = make(map[uint64]uint64)
	var buildingNameMap map[int64]string = make(map[int64]string)
	devices, err := d.deviceDao.GetByIds(c, deviceIds)
	if err != nil {
		zap.L().Error("读取设备列表失败", zap.Error(err))
		return deviceUtilizationList, errors.New("读取设备列表失败")
	}
	// 渲染建筑物
	for _, v := range devices {
		if v.DeviceBuildingId > 0 {
			buildingMap[v.DeviceId] = v.DeviceBuildingId
		}
	}
	var buildIds []uint64 = make([]uint64, len(buildingMap))
	buildIndex := 0
	for _, v := range buildingMap {
		buildIds[buildIndex] = v
		buildIndex++
	}
	// 读取建筑物列表
	buildList, err := d.deviceBuildDao.GetByIds(c, buildIds)
	if err != nil {
		zap.L().Error("读取建筑物列表失败", zap.Error(err))
		return deviceUtilizationList, err
	}

	for _, v := range buildList {
		buildingNameMap[v.ID] = v.Name
	}

	// 渲染设备列表
	for _, v := range devices {
		dId := int64(v.DeviceId)

		deviceMap[dId] = v
		deviceName := ""
		if v.Name != nil {
			deviceName = *v.Name
		}

		// 运行时间统计
		runTime := 0
		runTimeStr := "00"
		runRate := 0.0
		runInfo, ok := deviceRunMap[dId]
		if ok {
			runTime = int(runInfo.Value)
			if runTime > 0 {
				runTimeStr = time.SecondsToHMS(int64(runTime))
				runRate = math.RoundFloat(float64(float64(runInfo.Value)/float64(shiftTime))*100, 2)
			}
		}

		// 待机时间统计
		waitTime := 0
		waitTimeStr := "00"
		waitRate := 0.0
		waitInfo, ok := deviceWaitMap[dId]
		if ok {
			waitTime = int(waitInfo.Value)
			if waitTime > 0 {
				waitTimeStr = time.SecondsToHMS(int64(waitTime))
				waitRate = math.RoundFloat(float64(float64(waitInfo.Value)/float64(shiftTime))*100, 2)
			}
		}

		// 停机时间统计
		stopTime := 0
		stopTimeStr := "00"
		stopRate := 0.0
		if ok {
			stopTime = shiftTime - runTime - waitTime
			if stopTime < 0 {
				stopTime = 0
			}
			if stopTime > 0 {
				stopTimeStr = time.SecondsToHMS(int64(stopTime))
				stopRate = 100 - runRate - stopRate
			}
		}

		buildId := v.DeviceBuildingId
		buildName := ""
		buildName, ok = buildingNameMap[int64(buildId)]

		deviceUtilizationList = append(deviceUtilizationList, &deviceMonitorModel.DeviceUtilizationData{
			DeviceId:           dId,
			DeviceName:         deviceName,
			RunTime:            uint64(runTime),
			RunTimeStr:         runTimeStr,
			UtilizationRate:    runRate,
			UtilizationRateStr: fmt.Sprintf("%.2f", runRate) + "%",

			WaitTime:    uint64(waitTime),
			WaitTimeStr: waitTimeStr,
			WaitRate:    waitRate,
			WaitRateStr: fmt.Sprintf("%.2f", waitRate) + "%",

			StopTime:    uint64(stopTime),
			StopTimeStr: stopTimeStr,
			StopRate:    stopRate,
			StopRateStr: fmt.Sprintf("%.2f", stopRate) + "%",
			Building:    buildName,
		})
	}
	return deviceUtilizationList, nil
}
