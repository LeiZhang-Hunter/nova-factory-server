package device

import "fmt"

var (
	DEVICE_KEY = "DEVICE_KEY_%d"
)

func MakeDeviceKey(deviceId uint64) string {
	return fmt.Sprintf(DEVICE_KEY, deviceId)
}
