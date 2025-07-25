package device

import "fmt"

var (
	DEVICE_KEY = "DEVICE_KEY_%d"

	MODBUS_TCP = "modbus-tcp"
	MODBUS_RTU = "modbus-rtu"
	MLINK_TCP  = "mlink-tcp"
)

func MakeDeviceKey(deviceId uint64) string {
	return fmt.Sprintf(DEVICE_KEY, deviceId)
}

func CheckProtocol(protocol string) bool {
	if protocol != MODBUS_TCP && protocol != MODBUS_RTU && protocol != MLINK_TCP {
		return false
	}
	return true
}
