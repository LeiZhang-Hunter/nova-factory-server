package daemonizeserviceimpl

import (
	"fmt"
	"strconv"

	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	catalogresponse "nova-factory-server/app/business/iot/daemonize/daemonizemodels/response"
)

type catalogProtocolConverter interface {
	convert(*devicemodels.DeviceVO, *devicemodels.SysDeviceTemplate, []*devicemodels.SysModbusDeviceConfigData) (catalogresponse.Device, error)
}

// catalogConverterFor 根据设备协议选择目录转换器。
func catalogConverterFor(protocol string) (catalogProtocolConverter, error) {
	switch protocol {
	case "modbus-tcp", "modbus_tcp":
		return modbusTCPCatalogConverter{}, nil
	case "mqtt":
		return mqttCatalogConverter{}, nil
	default:
		return nil, fmt.Errorf("protocol %q is not supported by catalog", protocol)
	}
}

// catalogDevice 构造各协议共用的设备信息。
func catalogDevice(device *devicemodels.DeviceVO, template *devicemodels.SysDeviceTemplate, protocol string) catalogresponse.Device {
	item := catalogresponse.Device{ID: strconv.FormatUint(device.DeviceId, 10), Protocol: protocol, TemplateID: strconv.FormatInt(template.TemplateID, 10), TemplateName: template.Name, Points: []catalogresponse.Point{}}
	if device.Name != nil {
		item.Name = *device.Name
	}
	if device.Number != nil {
		item.Number = *device.Number
	}
	return item
}

// catalogPoint 构造各协议共用的测点信息。
func catalogPoint(point *devicemodels.SysModbusDeviceConfigData) catalogresponse.Point {
	return catalogresponse.Point{ID: strconv.FormatInt(point.DeviceConfigID, 10), Name: point.Name, Type: point.Type, Sort: point.Sort, Unit: point.Unit, Precision: point.Precision, Expression: point.Expression, DataFormat: point.DataFormat}
}
