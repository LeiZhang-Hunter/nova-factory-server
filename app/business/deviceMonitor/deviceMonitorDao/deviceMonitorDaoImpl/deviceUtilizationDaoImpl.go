package deviceMonitorDaoImpl

import (
	"errors"
	"fmt"
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
	"strconv"
	systime "time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DeviceUtilizationDaoImpl struct {
	iotDb          *iotdb.IotDb
	deviceDao      deviceDao.IDeviceDao
	deviceBuildDao buildingDao.BuildingDao
	shiftDao       systemDao.ISysShiftDao
	metricDao      metricDao.IMetricDao
	dictDataDao    systemDao.IDictDataDao
}

func NewDeviceUtilizationDaoImpl(iotDb *iotdb.IotDb, shiftDao systemDao.ISysShiftDao,
	deviceDao deviceDao.IDeviceDao,
	deviceBuildDao buildingDao.BuildingDao,
	metricDao metricDao.IMetricDao,
	dictDataDao systemDao.IDictDataDao) deviceMonitorDao.DeviceUtilizationDao {
	return &DeviceUtilizationDaoImpl{
		iotDb:          iotDb,
		shiftDao:       shiftDao,
		deviceDao:      deviceDao,
		deviceBuildDao: deviceBuildDao,
		metricDao:      metricDao,
		dictDataDao:    dictDataDao,
	}
}

// statRun 统计运行时设备
func (d *DeviceUtilizationDaoImpl) statDeviceStat(c *gin.Context, startTime string, endTime string,
	status int) (*deviceMonitorModel.DeviceStatusList, error) {
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
func (d *DeviceUtilizationDaoImpl) statDeviceProcess(c *gin.Context, startTime string, endTime string, interval string,
	status int) (*deviceMonitorModel.DeviceProcessList, error) {
	var processList deviceMonitorModel.DeviceProcessList
	processList.List = make(map[string][]deviceMonitorModel.DeviceStatus)
	session, err := d.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return nil, err
	}
	defer d.iotDb.PutSession(session)

	var timeout int64 = 5000

	sql := fmt.Sprintf("select sum(duration) as value from root.run_status_device.** where status = %d group by ([%s, %s), %s)  align by device",
		status, startTime, endTime, interval)

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
			Status: int((status)),
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

	var startTime string
	if req.Start > 0 {
		startTime = timeUtil.GetStartTime(req.Start, 200)
	} else {
		start := systime.Now().UnixMilli() - 86400*1000
		startTime = timeUtil.GetStartTime(uint64(start), 200)
	}
	endTime := timeUtil.GetEndTimeUseNow(req.End, true)

	// 查询所有班次，用来计算工作时间的稼动率
	// 使用Parse函数解析字符串为time.Time类型
	endTimeUnix, err := systime.Parse("2006-01-02 15:04:05", endTime)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return deviceUtilizationList, err
	}

	startTimeUnix, err := systime.Parse("2006-01-02 15:04:05", startTime)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return deviceUtilizationList, err
	}

	shiftTime := 0

	shiftTime = int(endTimeUnix.Unix() - startTimeUnix.Unix())

	if shiftTime <= 0 {
		shiftTime = 86400
	}

	// 从字典里读取全部状态
	runStatuses := d.dictDataDao.SelectDictDataByType(c, "device_run_status")

	// 批量读取设备id
	var deviceStatusMap map[int]map[int64]*deviceMonitorModel.DeviceStatus = make(map[int]map[int64]*deviceMonitorModel.DeviceStatus)
	var deviceMap map[int64]*deviceModels.DeviceVO = make(map[int64]*deviceModels.DeviceVO)

	for _, v := range runStatuses {
		statusValue, err := strconv.Atoi(v.DictValue)
		if err != nil {
			zap.L().Error("strconv.Atoi error", zap.Error(err))
			continue
		}
		statusList, err := d.statDeviceStat(c, startTime, endTime, statusValue)
		if err != nil {
			zap.L().Error("statDeviceStat error", zap.Error(err))
			continue
		}

		if statusList == nil {
			continue
		}

		for _, statusV := range statusList.List {
			if statusV.DeviceId == 0 {
				continue
			}
			deviceMap[statusV.DeviceId] = nil
			_, ok := deviceStatusMap[statusValue]
			if !ok {
				deviceStatusMap[statusValue] = make(map[int64]*deviceMonitorModel.DeviceStatus)
			}

			// Copy statusV to avoid appending the reference of loop variable
			statusValCopy := statusV
			deviceStatusMap[statusValue][statusValCopy.DeviceId] = &statusValCopy
		}
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

		// 动态统计各状态时间
		statusMap := make(map[int]deviceMonitorModel.DeviceStatusData)
		for status, statusValueMap := range deviceStatusMap {
			time := 0
			timeStr := "00"
			rate := 0.0
			runInfo, ok := statusValueMap[dId]
			if ok {
				time = int(runInfo.Value)
				if time > 0 {
					timeStr = timeUtil.SecondsToHMS(int64(time))
					rate = math.RoundFloat(float64(float64(runInfo.Value)/float64(shiftTime))*100, 2)
				}
			}
			statusMap[status] = deviceMonitorModel.DeviceStatusData{
				Time:    uint64(time),
				TimeStr: timeStr,
				Rate:    rate,
				RateStr: fmt.Sprintf("%.2f", rate) + "%",
			}
		}

		// 运行时间统计
		runTime := 0
		runTimeStr := "00"
		runRate := 0.0
		if runData, ok := statusMap[int(device.RUNNING)]; ok {
			runTime = int(runData.Time)
			runTimeStr = runData.TimeStr
			runRate = runData.Rate
		}

		// 待机时间统计
		waitTime := 0
		waitTimeStr := "00"
		waitRate := 0.0
		if waitData, ok := statusMap[int(device.WAITING)]; ok {
			waitTime = int(waitData.Time)
			waitTimeStr = waitData.TimeStr
			waitRate = waitData.Rate
		}

		// 停机时间统计
		stopTime := shiftTime - runTime - waitTime
		stopTimeStr := "00"
		stopRate := 0.0
		if stopTime < 0 {
			stopTime = 0
		}
		if stopTime > 0 {
			stopTimeStr = timeUtil.SecondsToHMS(int64(stopTime))
			stopRate = 100 - runRate - waitRate
		}

		buildId := v.DeviceBuildingId
		buildName := ""
		buildName, _ = buildingNameMap[int64(buildId)]

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
			StatusMap:   statusMap,
		})
	}
	return deviceUtilizationList, nil
}

// Search 统计稼动率
func (d *DeviceUtilizationDaoImpl) Search(c *gin.Context,
	req *deviceMonitorModel.DeviceUtilizationReq) (*deviceMonitorModel.DeviceUtilizationPublicDataList, error) {
	return &deviceMonitorModel.DeviceUtilizationPublicDataList{}, nil
}

// SearchV2 统计稼动率
func (d *DeviceUtilizationDaoImpl) SearchV2(c *gin.Context,
	req *deviceMonitorModel.DeviceUtilizationReq) (*deviceMonitorModel.DeviceUtilizationPublicDataListV2, error) {
	var deviceUtilizationList *deviceMonitorModel.DeviceUtilizationPublicDataListV2 = &deviceMonitorModel.DeviceUtilizationPublicDataListV2{
		List:        make([]*deviceMonitorModel.DeviceUtilizationDataV2, 0),
		ProcessList: make(map[string][]deviceMonitorModel.DeviceRunProcess, 0),
		Top5:        make(map[int][]*deviceMonitorModel.DeviceUtilizationDataV2),
	}

	if req == nil {
		return deviceUtilizationList, errors.New("参数错误")
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

	shiftTime := 86400

	// 从字典里读取全部状态
	runStatuses := d.dictDataDao.SelectDictDataByType(c, "device_run_status")

	// 批量读取设备id
	var deviceStatusMap map[int]map[int64]*deviceMonitorModel.DeviceStatus = make(map[int]map[int64]*deviceMonitorModel.DeviceStatus)
	var deviceMap map[int64]*deviceModels.DeviceVO = make(map[int64]*deviceModels.DeviceVO)

	for _, v := range runStatuses {
		statusValue, err := strconv.Atoi(v.DictValue)
		if err != nil {
			zap.L().Error("strconv.Atoi error", zap.Error(err))
			continue
		}
		statusList, err := d.statDeviceStat(c, startTime, endTime, statusValue)
		if err != nil {
			zap.L().Error("statDeviceStat error", zap.Error(err))
			continue
		}

		if statusList == nil {
			continue
		}

		for _, statusV := range statusList.List {
			if statusV.DeviceId == 0 {
				continue
			}
			deviceMap[statusV.DeviceId] = nil
			_, ok := deviceStatusMap[statusValue]
			if !ok {
				deviceStatusMap[statusValue] = make(map[int64]*deviceMonitorModel.DeviceStatus)
			}
			deviceStatusMap[statusValue][statusV.DeviceId] = &statusV
		}

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

	var deviceRunRateMap map[uint64]*deviceMonitorModel.DeviceUtilizationDataV2 = make(map[uint64]*deviceMonitorModel.DeviceUtilizationDataV2)

	// 渲染设备列表
	for _, v := range devices {
		dId := int64(v.DeviceId)

		deviceMap[dId] = v
		deviceName := ""
		if v.Name != nil {
			deviceName = *v.Name
		}

		data := &deviceMonitorModel.DeviceUtilizationDataV2{
			DeviceId:   dId,
			DeviceName: deviceName,
			StatusMap:  make(map[int]deviceMonitorModel.DeviceStatusData),
			Building:   "",
		}

		for status, statusValueMap := range deviceStatusMap {
			// 运行时间统计
			time := 0
			timeStr := "00"
			rate := 0.0
			runInfo, ok := statusValueMap[dId]
			if ok {
				time = int(runInfo.Value)
				if time > 0 {
					timeStr = timeUtil.SecondsToHMS(int64(time))
					rate = math.RoundFloat(float64(float64(runInfo.Value)/float64(shiftTime))*100, 2)
				}
			}
			if rate > 10 {
				fmt.Println(11)
			}
			data.StatusMap[status] = deviceMonitorModel.DeviceStatusData{
				Time:    uint64(time),
				TimeStr: timeStr,
				Rate:    rate,
				RateStr: fmt.Sprintf("%.2f", rate) + "%",
			}
		}

		buildId := v.DeviceBuildingId
		buildName := ""
		buildName, _ = buildingNameMap[int64(buildId)]

		data.Building = buildName
		deviceUtilizationList.List = append(deviceUtilizationList.List, data)
		deviceRunRateMap[uint64(dId)] = data
	}

	// processMapList 进度列表 status => DeviceProcessList
	processMapList := make(map[int]*deviceMonitorModel.DeviceProcessList)
	// 计算设备运行状态，间隔2小时一次]
	for _, status := range runStatuses {
		statusValue, err := strconv.Atoi(status.DictValue)
		if err != nil {
			zap.L().Error("strconv.Atoi error", zap.Error(err))
			continue
		}

		// 查询具体 设备运行情况，两个小时一个步长
		processList, err := d.statDeviceProcess(c, startTime, endTime, "2h", statusValue)
		if err != nil {
			zap.L().Error("statDeviceProcess error", zap.Error(err))
			continue
		}
		processMapList[statusValue] = processList
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

		if len(processMapList) == 0 {
			continue
		}

		// 设备运行进度
		var runProcess deviceMonitorModel.DeviceRunProcess
		if deviceValue.Name == nil {
			runProcess.DeviceName = ""
		} else {
			runProcess.DeviceName = *deviceValue.Name
		}
		runProcess.BuildingName = build

		deviceRunRateValue, ok := deviceRunRateMap[deviceValue.DeviceId]
		if ok {
			runProcess.StatusMap = deviceRunRateValue.StatusMap
		}

		for status, v := range processMapList {
			processValue, ok := v.List[deviceKey]
			if !ok {
				continue
			}

			for k, value := range processValue {
				if k >= len(runProcess.List) {
					statusValue := deviceMonitorModel.DeviceProcessStatus{
						Time:     value.Time,
						DeviceId: deviceId,
						Value:    make(map[int]float64),
					}
					statusValue.Value[status] = value.Value
					runProcess.List = append(runProcess.List, statusValue)
				} else {
					runProcess.List[k].Value[status] = value.Value
				}

			}
		}

		deviceUtilizationList.ProcessList[build] = append(deviceUtilizationList.ProcessList[build], runProcess)
	}

	// 计算设备概览
	statList, err := d.statDeviceRunStat(c, startTime, endTime)
	if err != nil {
		return nil, err
	}
	deviceUtilizationList.Total = uint64(len(statList))

	deviceUtilizationList.StatusCount = make(map[int]uint64)

	for _, stat := range runStatuses {
		statusValue, err := strconv.Atoi(stat.DictValue)
		if err != nil {
			zap.L().Error("strconv.Atoi error", zap.Error(err))
			continue
		}
		deviceUtilizationList.StatusCount[statusValue] = 0
	}

	// 统计各个状态的分值
	for _, stat := range statList {
		deviceUtilizationList.StatusCount[stat.Status]++
	}

	// 计算所有车间设备状态的top5
	for _, v := range runStatuses {
		statusValue, err := strconv.Atoi(v.DictValue)
		if err != nil {
			zap.L().Error("strconv.Atoi error", zap.Error(err))
			continue
		}
		statusSort := newDeviceStatusSort(statusValue, &deviceUtilizationList.List)
		sort.Sort(sort.Reverse(statusSort))
		for k, v := range deviceUtilizationList.List {
			if k >= 5 {
				break
			}
			deviceUtilizationList.Top5[statusValue] = append(deviceUtilizationList.Top5[statusValue], v)
		}
	}

	// 计算车间概括
	var totalTime map[string]uint64 = make(map[string]uint64)
	var statusCountTime map[string]map[int]uint64 = make(map[string]map[int]uint64)
	for _, v := range deviceUtilizationList.List {
		_, ok := totalTime[v.Building]
		if !ok {
			totalTime[v.Building] = 0
		}

		_, ok = statusCountTime[v.Building]
		if !ok {
			statusCountTime[v.Building] = make(map[int]uint64)
		}

		for status, count := range v.StatusMap {
			totalTime[v.Building] += count.Time
			statusCountTime[v.Building][status] += count.Time
		}
	}
	// 计算饼图占比
	// math.RoundFloat(float64(float64(waitInfo.Value)/float64(shiftTime))*100, 2)
	deviceUtilizationList.Radio = make(map[string][]deviceMonitorModel.DeviceRadio)
	for buildName, totalValue := range totalTime {
		_, ok := deviceUtilizationList.Radio[buildName]
		if !ok {
			deviceUtilizationList.Radio[buildName] = make([]deviceMonitorModel.DeviceRadio, 0)
		}

		deviceUtilizationList.Radio[buildName] = make([]deviceMonitorModel.DeviceRadio, 0)
		_, ok = statusCountTime[buildName]
		if !ok {
			continue
		}
		for status, count := range statusCountTime[buildName] {
			if totalValue != 0 {
				deviceUtilizationList.Radio[buildName] = append(deviceUtilizationList.Radio[buildName], deviceMonitorModel.DeviceRadio{
					Status:  status,
					Percent: math.RoundFloat(float64(float64(count)/float64(totalValue))*100, 2),
				})
			} else {
				deviceUtilizationList.Radio[buildName] = append(deviceUtilizationList.Radio[buildName], deviceMonitorModel.DeviceRadio{
					Status:  status,
					Percent: 0.00,
				})
			}

			continue

		}
	}

	//// 计算时间趋势
	end := systime.Now().UnixMilli()
	start := end - 7*86400*1000
	query := &metricModels.MetricQueryData{
		Values: make([]metricModels.MetricQueryValue, 0),
	}
	startTime = timeUtil.FormatDateFromSecond(start)
	endTime = timeUtil.FormatDateFromSecond(end)

	// 查询近一个周的运行情况
	processRunningList, err := d.statDeviceProcess(c, startTime, endTime, "1d", int(device.RUNNING))
	if err != nil {
		return nil, err
	}

	var processRunningListMap map[int64]float64 = make(map[int64]float64)
	if processRunningList != nil && len(processRunningList.List) > 0 {
		for _, v := range processRunningList.List {
			for _, runValue := range v {
				_, ok := processRunningListMap[runValue.Time]
				if !ok {
					processRunningListMap[runValue.Time] = 0
				}
				processRunningListMap[runValue.Time] += runValue.Value
			}
		}
	}

	// 1. 获取所有键并排序
	for k, v := range processRunningListMap {
		query.Values = append(query.Values, metricModels.MetricQueryValue{
			Time:  k,
			Value: v,
		})
	}

	// 按键升序排序
	sort.Slice(query.Values, func(i, j int) bool {
		return query.Values[i].Time < query.Values[j].Time
	})

	deviceUtilizationList.Data = query

	statusSort := newDeviceStatusSort(int(device.RUNNING), &deviceUtilizationList.List)
	sort.Sort(sort.Reverse(statusSort))
	return deviceUtilizationList, nil
}
