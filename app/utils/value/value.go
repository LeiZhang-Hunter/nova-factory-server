package value

type Type string

var (
	UNSPECIFIED Type = "unspecified"
	BIT         Type = "bit"
	UINT8       Type = "uint8"
	UINT16      Type = "uint16"
	UINT32      Type = "uint32"
	UINT64      Type = "uint64"
	INT8        Type = "int8"
	INT16       Type = "int16"
	INT32       Type = "int32"
	INT64       Type = "int64"
	FLOAT32     Type = "float32"
	FLOAT64     Type = "float64"
	STRING      Type = "string"
)

type Value struct {
	ByteValue   bool
	stringValue string
	boolValue   bool
	intValue    int64
	uintValue   uint64
	floatValue  float32
	doubleValue float64
}

type Data struct {
	Type  Type  `json:"type"`
	Value Value `json:"value"`
}
