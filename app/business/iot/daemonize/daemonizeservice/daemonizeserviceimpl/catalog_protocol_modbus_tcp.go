package daemonizeserviceimpl

import (
	"encoding/json"
	"fmt"

	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	catalogresponse "nova-factory-server/app/business/iot/daemonize/daemonizemodels/response"
)

type modbusTCPCatalogConverter struct{}

// convert 转换 Modbus TCP 连接和测点配置。
func (modbusTCPCatalogConverter) convert(device *devicemodels.DeviceVO, template *devicemodels.SysDeviceTemplate, points []*devicemodels.SysModbusDeviceConfigData) (catalogresponse.Device, error) {
	item := catalogDevice(device, template, "modbus_tcp")
	var extension devicemodels.ExtensionInfo
	if json.Unmarshal([]byte(device.Extension), &extension) != nil || len(extension.LocalInfo) != 1 {
		return item, fmt.Errorf("requires exactly one local connection")
	}
	connection := extension.LocalInfo[0]
	if connection.Slave <= 0 || connection.Quantity == 0 || connection.Address == "" {
		return item, fmt.Errorf("local connection is incomplete")
	}

	item.Config, _ = json.Marshal(map[string]any{"ipAddress": connection.Address, "slaveAddress": connection.Slave, "maxRegisters": connection.Quantity})
	for _, point := range points {
		if point.Register == nil || point.FunctionCode < 1 || point.FunctionCode > 4 {
			return item, fmt.Errorf("point %d has unsupported modbus measurement configuration", point.DeviceConfigID)
		}
		value := catalogPoint(point)
		value.Register = *point.Register
		value.FunctionCode = int8(point.FunctionCode)
		item.Points = append(item.Points, value)
	}
	return item, nil
}
