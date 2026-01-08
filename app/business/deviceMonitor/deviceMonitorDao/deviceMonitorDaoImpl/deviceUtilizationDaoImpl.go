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
	"nova-factory-server/app/business/metric/device/metricDao"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/business/system/systemDao"
	"nova-factory-server/app/constant/device"
	iotdb2 "nova-factory-server/app/constant/iotdb"
	"nova-factory-server/app/datasource/iotdb"
	"nova-factory-server/app/utils/math"
	timeUtil "nova-factory-server/app/utils/time"
	"sort"
	systime "time"
)

type DeviceUtilizationDaoImpl struct {
	iotDb          *iotdb.IotDb
	deviceDao      deviceDao.IDeviceDao
	deviceBuildDao buildingDao.BuildingDao
	shiftDao       systemDao.ISysShiftDao
	metricDao      metricDao.IMetricDao
}

func NewDeviceUtilizationDaoImpl(iotDb *iotdb.IotDb, shiftDao systemDao.ISysShiftDao,
	deviceDao deviceDao.IDeviceDao,
	deviceBuildDao buildingDao.BuildingDao,
	metricDao metricDao.IMetricDao) deviceMonitorDao.DeviceUtilizationDao {
	return &DeviceUtilizationDaoImpl{
		iotDb:          iotDb,
		shiftDao:       shiftDao,
		deviceDao:      deviceDao,
		deviceBuildDao: deviceBuildDao,
		metricDao:      metricDao,
	}
}

// statRun 统计运行时设备
func (d *DeviceUtilizationDaoImpl) statDeviceStat(c *gin.Context, startTime string, endTime string,
	status device.RUN_STATUS) (*deviceMonitorModel.DeviceStatusList, error) {
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

// statDeviceProcess 统计设备运行过程
func (d *DeviceUtilizationDaoImpl) statDeviceProcess(c *gin.Context, startTime string, endTime string,
	status device.RUN_STATUS) (*deviceMonitorModel.DeviceProcessList, error) {
	var processList deviceMonitorModel.DeviceProcessList
	processList.List = make(map[string][]deviceMonitorModel.DeviceStatus)
	session, err := d.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return nil, err
	}
	defer d.iotDb.PutSession(session)

	var timeout int64 = 5000

	sql := fmt.Sprintf("select sum(duration) as value from root.run_status_device.** where status = %d group by ([%s, %s), 2h)  align by device",
		status, startTime, endTime)

	statement, err := session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		zap.L().Error("读取设备运行过程失败:", zap.Error(err))
		return nil, err
	}
	for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
		timestamp := statement.GetTimestamp()
		deviceName := statement.GetText(statement.GetColumnName(0))
		duration := statement.GetDouble(statement.GetColumnName(1))
		_, ok := processList.List[deviceName]
		if !ok {
			processList.List[deviceName] = make([]deviceMonitorModel.DeviceStatus, 0)
		}
		processList.List[deviceName] = append(processList.List[deviceName], deviceMonitorModel.DeviceStatus{
			Value:  duration,
			Status: status,
			Time:   timestamp,
		})
		continue
	}

	return &processList, nil
}

// statDeviceProcess 统计设备运行过程
func (d *DeviceUtilizationDaoImpl) statDeviceRunStat(c *gin.Context, startTime string,
	endTime string) ([]deviceMonitorModel.DeviceRunStat, error) {
	var runStatList []deviceMonitorModel.DeviceRunStat = make([]deviceMonitorModel.DeviceRunStat, 0)
	session, err := d.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return nil, err
	}
	defer d.iotDb.PutSession(session)

	var timeout int64 = 5000

	sql := fmt.Sprintf("select last_value(status) from root.run_status_device.** where time>=%s and time < %s  align by device;",
		startTime, endTime)

	statement, err := session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		zap.L().Error("读取设备运行过程失败:", zap.Error(err))
		return nil, err
	}
	for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
		timestamp := statement.GetTimestamp()
		deviceName := statement.GetText(statement.GetColumnName(0))
		status := statement.GetInt64(statement.GetColumnName(1))
		var stat deviceMonitorModel.DeviceRunStat = deviceMonitorModel.DeviceRunStat{
			Time:   timestamp,
			Status: device.RUN_STATUS(status),
			Dev:    deviceName,
		}
		runStatList = append(runStatList, stat)
		continue
	}

	return runStatList, nil
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
		startTime = timeUtil.GetStartTime(req.Start, 200)
	} else {
		start := systime.Now().UnixMilli() - 86400*1000
		startTime = timeUtil.GetStartTime(uint64(start), 200)
	}
	endTime := timeUtil.GetEndTimeUseNow(req.End, true)

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
				runTimeStr = timeUtil.SecondsToHMS(int64(runTime))
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
				waitTimeStr = timeUtil.SecondsToHMS(int64(waitTime))
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
				stopTimeStr = timeUtil.SecondsToHMS(int64(stopTime))
				stopRate = 100 - runRate - waitRate
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

// Search 统计稼动率
func (d *DeviceUtilizationDaoImpl) Search(c *gin.Context,
	req *deviceMonitorModel.DeviceUtilizationReq) (*deviceMonitorModel.DeviceUtilizationPublicDataList, error) {
	var deviceUtilizationList *deviceMonitorModel.DeviceUtilizationPublicDataList = &deviceMonitorModel.DeviceUtilizationPublicDataList{
		List:        make([]*deviceMonitorModel.DeviceUtilizationData, 0),
		ProcessList: make(map[string][]deviceMonitorModel.DeviceRunProcess, 0),
		WaitTop5:    make([]*deviceMonitorModel.DeviceUtilizationData, 0),
		RunTop5:     make([]*deviceMonitorModel.DeviceUtilizationData, 0),
	}

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
		startTime = timeUtil.GetStartTime(req.Start, 200)
	} else {
		now := systime.Now()                                                                     // 获取当前时间
		midnight := systime.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()) // 当天0点的时间
		millDiff := now.Sub(midnight).Milliseconds()
		start := systime.Now().UnixMilli() - millDiff
		startTime = timeUtil.GetStartTime(uint64(start), 200)
	}
	endTime := timeUtil.GetEndTimeUseNow(req.End, true)

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

	var deviceRunRateMap map[uint64]*deviceMonitorModel.DeviceUtilizationData = make(map[uint64]*deviceMonitorModel.DeviceUtilizationData)

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
				runTimeStr = timeUtil.SecondsToHMS(int64(runTime))
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
				waitTimeStr = timeUtil.SecondsToHMS(int64(waitTime))
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
				stopTimeStr = timeUtil.SecondsToHMS(int64(stopTime))
				stopRate = 100 - runRate - waitRate
			}
		}

		buildId := v.DeviceBuildingId
		buildName := ""
		buildName, ok = buildingNameMap[int64(buildId)]
		data := &deviceMonitorModel.DeviceUtilizationData{
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
		}
		deviceUtilizationList.List = append(deviceUtilizationList.List, data)
		deviceRunRateMap[uint64(dId)] = data
	}

	// 查询具体 设备运行情况，两个小时一个步长
	processWaitList, err := d.statDeviceProcess(c, startTime, endTime, device.WAITING)
	if err != nil {
		return nil, err
	}

	processRunningList, err := d.statDeviceProcess(c, startTime, endTime, device.RUNNING)
	if err != nil {
		return nil, err
	}

	// 计算设备履历
	for deviceId, deviceValue := range deviceMap {
		deviceKey := iotdb2.MakeRunDeviceTemplateName(deviceId)
		build, ok := buildingNameMap[int64(deviceValue.DeviceBuildingId)]
		if !ok {
			continue
		}
		_, ok = deviceUtilizationList.ProcessList[build]
		if !ok {
			deviceUtilizationList.ProcessList[build] = make([]deviceMonitorModel.DeviceRunProcess, 0)
		}

		if processWaitList == nil {
			continue
		}

		if processRunningList == nil {
			continue
		}

		processWaitValue, ok := processWaitList.List[deviceKey]
		if !ok {
			continue
		}
		var runProcess deviceMonitorModel.DeviceRunProcess
		if deviceValue.Name == nil {
			runProcess.DeviceName = ""
		} else {
			runProcess.DeviceName = *deviceValue.Name
		}
		runProcess.BuildingName = build
		deviceRunRateValue, ok := deviceRunRateMap[deviceValue.DeviceId]
		if ok {
			runProcess.UtilizationRate = deviceRunRateValue.UtilizationRate
			runProcess.WaitRate = deviceRunRateValue.WaitRate
		}

		for _, v := range processWaitValue {
			statusValue := deviceMonitorModel.DeviceProcessStatus{
				Time:     v.Time,
				DeviceId: deviceId,
				Value:    make(map[device.RUN_STATUS]float64),
			}
			statusValue.Value[device.WAITING] = v.Value
			runProcess.List = append(runProcess.List, statusValue)
		}

		processRunningValue, ok := processRunningList.List[deviceKey]
		if !ok {
			continue
		}
		for k, v := range processRunningValue {
			if k > len(runProcess.List) {
				continue
			}
			runProcess.List[k].Time = v.Time
			runProcess.List[k].DeviceId = deviceId
			runProcess.List[k].Value[device.RUNNING] = v.Value

			stopProcessTime := 7200 - runProcess.List[k].Value[device.RUNNING] - runProcess.List[k].Value[device.WAITING]
			runProcess.List[k].Value[device.STOPPING] = stopProcessTime
		}
		deviceUtilizationList.ProcessList[build] = append(deviceUtilizationList.ProcessList[build], runProcess)
	}

	// 计算设备概览
	statList, err := d.statDeviceRunStat(c, startTime, endTime)
	if err != nil {
		return nil, err
	}
	deviceUtilizationList.Total = uint64(len(statList))
	for _, stat := range statList {
		if stat.Status == device.RUNNING {
			deviceUtilizationList.Running++
		} else if stat.Status == device.WAITING {
			deviceUtilizationList.WaitTing++
		} else {
			deviceUtilizationList.Stopped++
		}
	}

	// 计算运行时间top5
	sort.Sort(sort.Reverse(runSortUtilization(deviceUtilizationList.List)))
	for k, v := range deviceUtilizationList.List {
		if k >= 5 {
			break
		}
		deviceUtilizationList.RunTop5 = append(deviceUtilizationList.RunTop5, v)
	}

	// 计算待机率top5
	sort.Sort(sort.Reverse(waitSortUtilization(deviceUtilizationList.List)))
	for k, v := range deviceUtilizationList.List {
		if k >= 5 {
			break
		}
		deviceUtilizationList.WaitTop5 = append(deviceUtilizationList.WaitTop5, v)
	}

	// 计算车间概括
	var totalTime map[string]uint64 = make(map[string]uint64)
	var runTime map[string]uint64 = make(map[string]uint64)
	var stopTime map[string]uint64 = make(map[string]uint64)
	var waitTime map[string]uint64 = make(map[string]uint64)
	for _, v := range deviceUtilizationList.List {
		_, ok := totalTime[v.Building]
		if !ok {
			totalTime[v.Building] = 0
		}

		_, ok = runTime[v.Building]
		if !ok {
			runTime[v.Building] = 0
		}

		_, ok = stopTime[v.Building]
		if !ok {
			stopTime[v.Building] = 0
		}

		_, ok = waitTime[v.Building]
		if !ok {
			waitTime[v.Building] = 0
		}
		runTime[v.Building] += v.RunTime
		stopTime[v.Building] += v.StopTime
		waitTime[v.Building] += v.WaitTime
		totalTime[v.Building] += v.WaitTime + v.StopTime + v.RunTime
	}
	// math.RoundFloat(float64(float64(waitInfo.Value)/float64(shiftTime))*100, 2)
	deviceUtilizationList.Radio = make(map[string][]deviceMonitorModel.DeviceRadio)
	for buildName, v := range totalTime {
		_, ok := deviceUtilizationList.Radio[buildName]
		if !ok {
			deviceUtilizationList.Radio[buildName] = make([]deviceMonitorModel.DeviceRadio, 0)
		}

		runValue, ok := runTime[buildName]
		stopValue, ok := stopTime[buildName]
		waitValue, ok := waitTime[buildName]
		deviceUtilizationList.Radio[buildName] = []deviceMonitorModel.DeviceRadio{
			{
				Status:  device.RUNNING,
				Percent: math.RoundFloat(float64(float64(runValue)/float64(v))*100, 2),
			},
			{
				Status:  device.WAITING,
				Percent: math.RoundFloat(float64(float64(waitValue)/float64(v))*100, 2),
			},
			{
				Status:  device.STOPPING,
				Percent: math.RoundFloat(float64(float64(stopValue)/float64(v))*100, 2),
			},
		}
	}

	// 计算时间趋势
	end := systime.Now().UnixMilli()
	start := end - 7*86400*1000
	query, err := d.metricDao.Query(c, &metricModels.MetricDataQueryReq{
		Type:       "line",
		Name:       "root.run_status_device.**",
		Start:      uint64(start),
		Field:      " ",
		End:        uint64(end),
		Step:       0,
		Interval:   1440,
		Expression: "Sum(duration)  as value",
		Having:     fmt.Sprintf("last_value(status) = %d", device.WAITING),
	})
	if err != nil {
		return deviceUtilizationList, err
	}
	deviceUtilizationList.Data = query
	sort.Sort(sort.Reverse(runSortUtilization(deviceUtilizationList.List)))
	return deviceUtilizationList, nil
}
