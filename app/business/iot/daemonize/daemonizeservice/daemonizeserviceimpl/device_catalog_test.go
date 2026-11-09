package daemonizeserviceimpl

import (
	"encoding/json"
	"strings"
	"testing"

	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	catalogresponse "nova-factory-server/app/business/iot/daemonize/daemonizemodels/response"
)

// TestBuildDeviceCatalogMixedProtocols 验证混合协议目录及各协议专属校验。
func TestBuildDeviceCatalogMixedProtocols(t *testing.T) {
	register := int64(42)
	snapshot := &catalogresponse.CatalogSnapshot{
		Devices: []*devicemodels.DeviceVO{
			{DeviceId: 1, DeviceProtocolId: 10, ProtocolType: "modbus-tcp", Extension: `{"localInfo":[{"address":"host:502","slave":1,"quantity":10}]}`},
			{DeviceId: 2, DeviceProtocolId: 20, ProtocolType: "mqtt", Extension: `{"mqttInfo":[{"address":"tcp://broker:1883","topic":"sensors","client_id":"client"}]}`},
		},
		Templates: []*devicemodels.SysDeviceTemplate{{TemplateID: 10}, {TemplateID: 20}},
		Points: []*devicemodels.SysModbusDeviceConfigData{
			{DeviceConfigID: 101, TemplateID: 10, Register: &register, FunctionCode: 3},
			{DeviceConfigID: 201, TemplateID: 20, Annotation: `[{"key":"expression","value":"temperature"}]`},
		},
	}

	catalog, err := buildDeviceCatalog(9, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if len(catalog.Devices) != 2 || catalog.Devices[0].Points[0].Register != 42 || catalog.Devices[1].Points[0].Expression != "temperature" || catalog.Devices[1].Points[0].DataFormat != "float64" || catalog.Version == "" {
		t.Fatalf("unexpected catalog: %+v", catalog)
	}
	var config devicemodels.MqttInfo
	if err := json.Unmarshal(catalog.Devices[1].Config, &config); err != nil || config.Topic != "sensors" || config.ClientId != "client" {
		t.Fatalf("unexpected MQTT config: %+v, %v", config, err)
	}

	snapshot.Points[0].FunctionCode = 0
	if _, err := buildDeviceCatalog(9, snapshot); err == nil || !strings.Contains(err.Error(), "point 101") {
		t.Fatalf("expected Modbus point error, got %v", err)
	}

	snapshot.Devices[0].ProtocolType = "modbus-rtu"
	if _, err := buildDeviceCatalog(9, snapshot); err == nil || !strings.Contains(err.Error(), "modbus-rtu") {
		t.Fatalf("expected unsupported protocol error, got %v", err)
	}
}
