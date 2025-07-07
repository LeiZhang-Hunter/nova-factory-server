package metricDaoIMpl

import (
	"context"
	"fmt"
	"github.com/apache/iotdb-client-go/client"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/metric/device/metricModels"
	iotdb2 "nova-factory-server/app/constant/iotdb"
	"nova-factory-server/app/datasource/iotdb"
	"nova-factory-server/app/utils/time"
)

type iotDbExport struct {
	iotDb *iotdb.IotDb
}

func newIotDbExport(iotDb *iotdb.IotDb) iDaoExport {
	return &iotDbExport{
		iotDb: iotDb,
	}
}

func (i *iotDbExport) Export(ctx context.Context, data []*metricModels.NovaMetricsDevice) error {
	session, err := i.iotDb.GetSession()
	if err != nil {
		return err
	}
	defer i.iotDb.PutSession(session)

	// 变更数据结构，根据设备id分类
	var dataMap map[uint64][]*metricModels.NovaMetricsDevice = make(map[uint64][]*metricModels.NovaMetricsDevice)
	for _, d := range data {
		_, ok := dataMap[d.DeviceId]
		if !ok {
			dataMap[d.DeviceId] = make([]*metricModels.NovaMetricsDevice, 0)
		}
		dataMap[d.DeviceId] = append(dataMap[d.DeviceId], d)
	}

	for deviceId, list := range dataMap {
		name := fmt.Sprintf(iotdb2.ROOT_DEVICE_TEMPLATE_NAME, deviceId)
		measurementsSlice := [][]string{}
		dataTypes := [][]client.TSDataType{}
		values := [][]interface{}{}
		timestamps := []int64{}
		for _, value := range list {
			measurementsSlice = append(measurementsSlice, []string{
				"template_id", "data_id", "value",
			})
			dataTypes = append(dataTypes, []client.TSDataType{client.INT64, client.INT64, client.FLOAT})
			values = append(values, []interface{}{int64(value.TemplateId), int64(value.DataId), float32(value.Value)})
			timestamps = append(timestamps, value.StartTimeUnix.UnixMilli())
		}
		stat, err := session.InsertRecordsOfOneDevice(name, timestamps, measurementsSlice, dataTypes, values, false)
		if err != nil {
			zap.L().Error("InsertRecordsOfOneDevice error", zap.Error(err))
			continue
		}
		fmt.Println(stat)
	}

	return nil
}

func (i *iotDbExport) Metric(c *gin.Context, req *metricModels.MetricQueryReq) (*metricModels.MetricQueryData, error) {
	if req == nil {
		return nil, nil
	}

	session, err := i.iotDb.GetSession()
	if err != nil {
		return nil, err
	}
	defer i.iotDb.PutSession(session)

	var startTime string
	if req.Start > 0 {
		startTime = time.GetStartTime(req.Start, 200)
	}
	endTime := time.GetEndTimeUseNow(req.End, true)

	whereSql := "where "
	if startTime != "" && endTime != "" {
		whereSql += "Time >= " + startTime
	}
	if endTime != "" {
		whereSql += "and Time < " + endTime
	}
	if req.Step <= 0 {
		req.Step = 1
	}
	return nil, nil

}

// InstallDevice 安装设备模板
func (i *iotDbExport) InstallDevice(c *gin.Context, device *deviceModels.DeviceVO) error {
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return err
	}
	defer i.iotDb.PutSession(session)

	name := fmt.Sprintf(iotdb2.ROOT_DEVICE_TEMPLATE_NAME, device.DeviceId)
	// 创建设备模板
	group, err := session.SetStorageGroup(name)
	if err != nil {
		zap.L().Error("创建设备数据库失败, ", zap.Error(err), zap.Any("code", group.GetCode()))
		return err
	}

	// 挂载设备模板
	_, err = session.ExecuteStatement(fmt.Sprintf("set device template %s to %s", iotdb2.NOVA_DEVICE_TEMPLATE, name))
	if err != nil {
		zap.L().Error("绑定设备数据库失败, ", zap.Error(err))
		return err
	}

	// 激活设备模板
	_, err = session.ExecuteStatement(fmt.Sprintf("create timeseries using device template on %s", name))
	if err != nil {
		zap.L().Error("激活设备模板失败, ", zap.Error(err))
		return err
	}
	return nil
}

// UnInStallDevice 卸载设备模板
func (i *iotDbExport) UnInStallDevice(c *gin.Context, deviceId int64) error {
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return err
	}
	defer i.iotDb.PutSession(session)

	name := fmt.Sprintf(iotdb2.ROOT_DEVICE_TEMPLATE_NAME, deviceId)

	// 删除模板表示的某一组时间序列
	_, err = session.ExecuteStatement(fmt.Sprintf("deactivate device template %s from %s", iotdb2.NOVA_DEVICE_TEMPLATE, name))
	if err != nil {
		zap.L().Error("deactivate  device template", zap.Error(err))
		return err
	}

	_, err = session.ExecuteStatement(fmt.Sprintf("unset device template %s from %s", iotdb2.NOVA_DEVICE_TEMPLATE, name))
	if err != nil {
		zap.L().Error("unset  device template", zap.Error(err))
		return err
	}

	_, err = session.ExecuteStatement(fmt.Sprintf("drop database %s", name))
	if err != nil {
		zap.L().Error("unset  device template", zap.Error(err))
		return err
	}
	return nil
}
