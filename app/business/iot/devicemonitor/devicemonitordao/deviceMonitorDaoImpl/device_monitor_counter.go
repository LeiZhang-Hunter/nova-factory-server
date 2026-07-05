package deviceMonitorDaoImpl

import (
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitordao"
	deviceMonitorModel2 "nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/metric/device/metricdao"
	"nova-factory-server/app/business/iot/metric/device/metricmodels"
	"nova-factory-server/app/business/iot/metric/device/metricservice"

	"github.com/gin-gonic/gin"
)

type DeviceMonitorCalcDaoImpl struct {
	metricDao  metricdao.IMetricDao
	devService metricservice.IDevMapService
}

func NewDeviceMonitorCalcDaoImpl(metricDao metricdao.IMetricDao, devService metricservice.IDevMapService) devicemonitordao.DeviceMonitorCalcDao {
	return &DeviceMonitorCalcDaoImpl{
		metricDao:  metricDao,
		devService: devService,
	}
}

func (dao *DeviceMonitorCalcDaoImpl) CounterByTimeRange(startTime int64, endTime int64, interval string) (*metricmodels.MetricQueryData, error) {
	return dao.metricDao.CounterByTimeRange(startTime, endTime, interval)
}

func (dao *DeviceMonitorCalcDaoImpl) CounterByDevice(c *gin.Context, startTime int64, endTime int64, limit int) (*deviceMonitorModel2.TypeDeviceCounterRank, error) {
	rank, err := dao.metricDao.CounterByDevice(c, startTime, endTime, limit)
	if err != nil {
		return nil, err
	}

	if rank == nil {
		return &deviceMonitorModel2.TypeDeviceCounterRank{Rows: make([]*deviceMonitorModel2.TypeDeviceCounterRankValue, 0)}, nil
	}

	if rank.Rows == nil || len(rank.Rows) == 0 {
		return rank, nil
	}

	var devs []string = make([]string, 0)
	for _, value := range rank.Rows {
		devs = append(devs, value.Dev)
	}
	if len(devs) == 0 {
		return rank, nil
	}

	devMapList, err := dao.devService.GetDevList(c, devs)
	if err != nil {
		return nil, err
	}

	var devDataMap map[string]deviceMonitorModel2.SysIotDbDevMapData = make(map[string]deviceMonitorModel2.SysIotDbDevMapData)
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

	return rank, nil
}
