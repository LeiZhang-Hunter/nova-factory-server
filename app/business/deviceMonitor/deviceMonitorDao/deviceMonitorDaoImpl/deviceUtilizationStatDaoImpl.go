package deviceMonitorDaoImpl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	iotdb2 "nova-factory-server/app/constant/iotdb"
)

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

// statDeviceStatByDeviceId 统计运行时设备
func (d *DeviceUtilizationDaoImpl) statDeviceStatByDeviceId(c *gin.Context, startTime string, endTime string,
	deviceId int64, status int) (*deviceMonitorModel.DeviceStatusList, error) {
	session, err := d.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return nil, err
	}
	defer d.iotDb.PutSession(session)
	deviceKey := iotdb2.MakeRunDeviceTemplateName(deviceId)
	var timeout int64 = 5000

	sql := fmt.Sprintf("select sum(duration) as value from %s where time > %s and time < %s and status = %d align by device",
		deviceKey, startTime, endTime, status)

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

// statDeviceProcessByDeviceId 统计设备运行过程
func (d *DeviceUtilizationDaoImpl) statDeviceProcessByDeviceId(c *gin.Context, startTime string, endTime string,
	deviceId int64, interval string,
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

	deviceKey := iotdb2.MakeRunDeviceTemplateName(deviceId)

	sql := fmt.Sprintf("select sum(duration) as value from %s where status = %d group by ([%s, %s), %s)  align by device",
		deviceKey, status, startTime, endTime, interval)

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
