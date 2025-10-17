package render

import (
	"encoding/json"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/constant/gateway"
	"nova-factory-server/app/utils/gateway/v1/api"
	"nova-factory-server/app/utils/gateway/v1/config/app/source/bhps7"
	"nova-factory-server/app/utils/gateway/v1/config/app/source/mqtt"
	"strconv"
)

// OfBhps7Device modbus tcp 设备数据
func OfBhps7Device(vo *deviceModels.DeviceVO, data []*deviceModels.SysModbusDeviceConfigData) *bhps7.Device {
	var device bhps7.Device
	device.DeviceId = strconv.FormatUint(vo.DeviceId, 10)
	if vo.Name != nil {
		device.Name = *vo.Name
	} else {
		device.Name = strconv.FormatUint(vo.DeviceId, 10)
	}
	device.Template = make([]bhps7.DataTypeConfig, 0)
	for _, v := range data {
		device.Template = append(device.Template, bhps7.DataTypeConfig{
			Name:         v.Name,
			Protocol:     "",
			DataId:       uint64(v.DeviceConfigID),
			Annotation:   v.Name,
			DataFormat:   api.DataValueType(v.DataFormat),
			Unit:         v.Unit,
			TemplateId:   uint64(v.TemplateID),
			Position:     uint16(v.Register),
			FunctionCode: uint16(v.FunctionCode),
			Sort:         api.ByteOrder(v.Sort),
		})
	}

	return &device
}

// OfMqttDevice mqtt tcp 设备数据
func OfMqttDevice(vo *deviceModels.DeviceVO, data []*deviceModels.SysModbusDeviceConfigData) *mqtt.Device {
	var device mqtt.Device
	device.DeviceId = strconv.FormatUint(vo.DeviceId, 10)
	if vo.Name != nil {
		device.Name = *vo.Name
	} else {
		device.Name = strconv.FormatUint(vo.DeviceId, 10)
	}
	device.Template = make([]mqtt.DataTypeConfig, 0)
	for _, v := range data {
		annotations := make([]deviceModels.Annotation, 0)
		err := json.Unmarshal([]byte(v.Annotation), &annotations)
		if err != nil {
			zap.L().Error("unmarshal annotations error", zap.Error(err))
			continue
		}
		var expression string
		for _, a := range annotations {
			if a.Key == gateway.Expression {
				expression = a.Value
			}
		}

		if v.DataFormat == "" {
			v.DataFormat = string(api.Float64)
		}

		device.Template = append(device.Template, mqtt.DataTypeConfig{
			Name:       v.Name,
			DataId:     uint64(v.DeviceConfigID),
			DataFormat: api.DataValueType(v.DataFormat),
			Unit:       v.Unit,
			TemplateId: uint64(v.TemplateID),
			Expression: expression,
		})
	}

	return &device
}
