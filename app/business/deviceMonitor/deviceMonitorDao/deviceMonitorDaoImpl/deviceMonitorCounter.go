package deviceMonitorDaoImpl

import (
	"fmt"
	"github.com/apache/iotdb-client-go/client"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorDao"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/business/metric/device/metricService"
	"nova-factory-server/app/datasource/iotdb"
	"nova-factory-server/app/utils/math"
	"nova-factory-server/app/utils/time"
)

type DeviceMonitorCalcDaoImpl struct {
	iotDb      *iotdb.IotDb
	devService metricService.IDevMapService
}

func NewDeviceMonitorCalcDaoImpl(iotDb *iotdb.IotDb, devService metricService.IDevMapService) deviceMonitorDao.DeviceMonitorCalcDao {
	return &DeviceMonitorCalcDaoImpl{
		iotDb:      iotDb,
		devService: devService,
	}
}

func (dao *DeviceMonitorCalcDaoImpl) CounterByTimeRange(startTime int64, endTime int64, interval string) (*metricModels.MetricQueryData, error) {

	session, err := dao.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return nil, err
	}
	defer dao.iotDb.PutSession(session)

	if interval == "" {
		intervalValue := (endTime - startTime) / 60 / 30 / 1000
		interval = fmt.Sprintf("%dm", intervalValue)
	}

	var timeout int64 = 5000
	var data *metricModels.MetricQueryData = metricModels.NewMetricQueryData()

	sql := fmt.Sprintf("select count(value) as value from root.device.** group by([%s, %s), %s),level=1",
		time.GetStartTime(uint64(startTime), 0), time.GetEndTime(uint64(endTime), 0), interval)

	statement, err := session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		zap.L().Error("ExecuteQueryStatement error", zap.Error(err))
		return nil, err
	}
	if len(statement.GetColumnNames()) <= 1 {
		for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
			timestamp := statement.GetTimestamp()
			var v float64
			dataType := statement.GetColumnDataType(0)
			switch dataType {
			case client.BOOLEAN:
				{
					dataValue := statement.GetBool(statement.GetColumnName(0))
					if dataValue == true {
						v = 1.0
					} else {
						v = 0.0
					}
					break
				}
			case client.INT32:
				{
					dataValue := statement.GetInt32(statement.GetColumnName(0))
					v = float64(dataValue)
					break
				}
			case client.INT64:
				{
					dataValue := statement.GetInt64(statement.GetColumnName(0))
					v = float64(dataValue)
					break
				}
			case client.FLOAT:
				{
					dataValue := statement.GetFloat(statement.GetColumnName(0))
					v = float64(dataValue)
					break
				}
			case client.DOUBLE:
				{
					dataValue := statement.GetDouble(statement.GetColumnName(0))
					v = float64(dataValue)
					break
				}

			}

			data.Values = append(data.Values, metricModels.MetricQueryValue{
				Time:  timestamp,
				Value: math.RoundFloat(v, 2),
			})

		}
	} else {
		data.MultiValues = make([][]metricModels.MetricQueryValue, len(statement.GetColumnNames()))

		for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
			timestamp := statement.GetTimestamp()

			for k, column := range statement.GetColumnNames() {
				var v float64
				dataType := statement.GetColumnDataType(0)
				switch dataType {
				case client.BOOLEAN:
					{
						dataValue := statement.GetBool(column)
						if dataValue == true {
							v = 1.0
						} else {
							v = 0.0
						}
						break
					}
				case client.INT32:
					{
						dataValue := statement.GetInt32(column)
						v = float64(dataValue)
						break
					}
				case client.INT64:
					{
						dataValue := statement.GetInt64(column)
						v = float64(dataValue)
						break
					}
				case client.FLOAT:
					{
						dataValue := statement.GetFloat(column)
						v = float64(dataValue)
						break
					}
				case client.DOUBLE:
					{
						dataValue := statement.GetDouble(column)
						v = float64(dataValue)
						break
					}

				}
				data.MultiValues[k] = append(data.MultiValues[k], metricModels.MetricQueryValue{
					Time:  timestamp,
					Value: math.RoundFloat(v, 2),
				})
			}
		}
	}

	return data, nil

}

func (dao *DeviceMonitorCalcDaoImpl) CounterByDevice(c *gin.Context, startTime int64, endTime int64, limit int) (*deviceMonitorModel.TypeDeviceCounterRank, error) {
	session, err := dao.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return nil, err
	}
	defer dao.iotDb.PutSession(session)

	var timeout int64 = 5000

	sql := fmt.Sprintf("select count(value) from root.device.** where time > %s and time < %s  order by count(value) desc limit %d ALIGN BY DEVICE",
		time.GetStartTime(uint64(startTime), 0), time.GetEndTime(uint64(endTime), 0), limit)

	statement, err := session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		zap.L().Error("ExecuteQueryStatement error", zap.Error(err))
		return nil, err
	}

	rank := deviceMonitorModel.TypeDeviceCounterRank{
		Rows: make([]*deviceMonitorModel.TypeDeviceCounterRankValue, 0),
	}

	for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
		timestamp := statement.GetTimestamp()
		device := statement.GetText(statement.GetColumnName(0))
		value := statement.GetInt64(statement.GetColumnName(1))
		rank.Rows = append(rank.Rows, &deviceMonitorModel.TypeDeviceCounterRankValue{
			Time:  timestamp,
			Dev:   device,
			Value: value,
		})
	}

	if rank.Rows == nil || len(rank.Rows) == 0 {
		return &rank, nil
	}

	var devs []string = make([]string, 0)
	for _, value := range rank.Rows {
		devs = append(devs, value.Dev)
	}
	if len(devs) == 0 {
		return &rank, nil
	}

	devMapList, err := dao.devService.GetDevList(c, devs)
	if err != nil {
		return nil, err
	}

	var devDataMap map[string]deviceMonitorModel.SysIotDbDevMapData = make(map[string]deviceMonitorModel.SysIotDbDevMapData)
	for _, v := range devMapList {
		devDataMap[v.Device] = v
	}

	for k, value := range rank.Rows {
		devValue, ok := devDataMap[value.Dev]
		if !ok {
			continue
		}
		rank.Rows[k].DevName = devValue.DevName
		rank.Rows[k].DataName = devValue.DataName
	}

	return &rank, nil
}
