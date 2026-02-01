package device

import "fmt"

type RUN_STATUS int

var STOPPING RUN_STATUS = 0
var WAITING RUN_STATUS = 1
var RUNNING RUN_STATUS = 2

var (
	DEVICE_KEY = "DEVICE_KEY_%d"

	DEVICE_CONTROL_KEY = "DEVICE_CONTROL_KEY_%d_%d" // device_id data_id

	MODBUS_TCP = "modbus-tcp"
	MODBUS_RTU = "modbus-rtu"
	MLINK_TCP  = "mlink-tcp"
	MQTT       = "mqtt"
)

func MakeDeviceKey(deviceId uint64) string {
	return fmt.Sprintf(DEVICE_KEY, deviceId)
}

func CheckProtocol(protocol string) bool {
	if protocol != MODBUS_TCP && protocol != MODBUS_RTU && protocol != MLINK_TCP && protocol != MQTT {
		return false
	}
	return true
}
