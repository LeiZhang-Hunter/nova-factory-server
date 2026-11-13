package iotdb

import (
	"testing"

	"github.com/apache/iotdb-client-go/client"
)

// TestIoTDBSampleValue 验证整数写入类型和大整数精度，以及浮点写入类型。
func TestIoTDBSampleValue(t *testing.T) {
	const largeInteger int64 = 9007199254740993
	integer := NewIntMetricSample("root.run_status_device.dev1.duration", nil, 1794590079000, largeInteger)
	dataType, value := iotDBSampleValue(integer)
	if dataType != client.INT64 || value != largeInteger {
		t.Fatalf("integer sample: got (%v, %v), want (INT64, %d)", dataType, value, largeInteger)
	}

	decimal := NewMetricSample("root.device.dev1.value", nil, 1794590079000, 60.5)
	dataType, value = iotDBSampleValue(decimal)
	if dataType != client.DOUBLE || value != 60.5 {
		t.Fatalf("double sample: got (%v, %v), want (DOUBLE, 60.5)", dataType, value)
	}
}
