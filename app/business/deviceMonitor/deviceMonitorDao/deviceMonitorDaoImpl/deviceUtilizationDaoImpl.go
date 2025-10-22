package deviceMonitorDaoImpl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorDao"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/constant/device"
	"nova-factory-server/app/datasource/iotdb"
	"nova-factory-server/app/utils/time"
	systime "time"
)

type DeviceUtilizationDaoImpl struct {
	iotDb *iotdb.IotDb
}

func NewDeviceUtilizationDaoImpl(iotDb *iotdb.IotDb) deviceMonitorDao.DeviceUtilizationDao {
	return &DeviceUtilizationDaoImpl{
		iotDb: iotDb,
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
		fmt.Println(v)
		fmt.Println(deviceName)
		var deviceId int64
		fmt.Sscanf(deviceName, "root.run_status_device.dev%d", &deviceId)
		data.List = append(data.List, deviceMonitorModel.DeviceStatus{
			DeviceId: deviceId,
			Value:    v,
			Status:   status,
		})
	}

	return data, nil
}

// Stat 统计稼动率
func (d *DeviceUtilizationDaoImpl) Stat(c *gin.Context, req *deviceMonitorModel.DeviceUtilizationReq) {
	if req == nil {
		return
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
		return
	}
	waitingList, err := d.statDeviceStat(c, startTime, endTime, device.WAITING)
	if err != nil {
		return
	}
	fmt.Println(waitingList)
	fmt.Println(runList)
	return
}
