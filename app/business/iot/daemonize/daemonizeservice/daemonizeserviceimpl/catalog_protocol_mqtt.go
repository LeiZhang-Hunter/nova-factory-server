package daemonizeserviceimpl

import (
	"encoding/json"
	"fmt"

	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	catalogresponse "nova-factory-server/app/business/iot/daemonize/daemonizemodels/response"
	"nova-factory-server/app/constant/gateway"
)

type mqttCatalogConverter struct{}

// convert 转换 MQTT 连接和测点配置。
func (mqttCatalogConverter) convert(device *devicemodels.DeviceVO, template *devicemodels.SysDeviceTemplate, points []*devicemodels.SysModbusDeviceConfigData) (catalogresponse.Device, error) {
	item := catalogDevice(device, template, "mqtt")
	var extension devicemodels.ExtensionInfo
	if json.Unmarshal([]byte(device.Extension), &extension) != nil || len(extension.LocalMqttInfo) != 1 {
		return item, fmt.Errorf("requires exactly one MQTT connection")
	}
	connection := extension.LocalMqttInfo[0]
	if connection.Address == "" || connection.Topic == "" {
		return item, fmt.Errorf("MQTT address or topic is missing")
	}

	item.Config, _ = json.Marshal(connection)
	for _, point := range points {
		value := catalogPoint(point)
		var annotations []devicemodels.Annotation
		if point.Annotation != "" {
			if err := json.Unmarshal([]byte(point.Annotation), &annotations); err != nil {
				return item, fmt.Errorf("point %d has invalid annotations: %w", point.DeviceConfigID, err)
			}
		}
		for _, annotation := range annotations {
			if annotation.Key == gateway.Expression {
				value.Expression = annotation.Value
			}
		}
		if value.DataFormat == "" {
			value.DataFormat = "float64"
		}
		item.Points = append(item.Points, value)
	}
	return item, nil
}
