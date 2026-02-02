package value

import (
	"errors"
	"fmt"
	controlService "github.com/novawatcher-io/nova-factory-payload/control/v1"
	"strconv"
)

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
	ByteValue   bool    `json:"byteValue"`
	StringValue string  `json:"stringValue"`
	BoolValue   bool    `json:"boolValue"`
	IntValue    int64   `json:"intValue"`
	UintValue   uint64  `json:"uintValue"`
	FloatValue  float32 `json:"floatValue"`
	DoubleValue float64 `json:"doubleValue"`
}

type Data struct {
	Type  Type  `json:"type"`
	Value Value `json:"value"`
}

func (d *Data) ToString() string {
	switch d.Type {
	case BIT:
		return strconv.FormatBool(d.Value.BoolValue)
	case UINT8, UINT16, UINT32, UINT64:
		return strconv.FormatUint(d.Value.UintValue, 10)
	case INT8, INT16, INT32, INT64:
		return strconv.FormatInt(d.Value.IntValue, 10)
	case FLOAT32:
		return fmt.Sprintf("%f", d.Value.FloatValue)
	case FLOAT64:
		return fmt.Sprintf("%f", d.Value.DoubleValue)
	case STRING:
		return d.Value.StringValue
	default:
		return ""
	}
}

func (d *Data) ToValue() (*controlService.Value, error) {
	var value controlService.Value
	tool := NewValueTool(&value)
	switch d.Type {
	case BIT:
		{
			err := tool.SetBit(d.Value)
			if err != nil {
				return nil, err
			}
			return &value, nil
		}
	case UINT8:
		{
			err := tool.SetUint(d.Value, controlService.ValueType_VALUE_TYPE_UINT8)
			if err != nil {
				return nil, err
			}
			return &value, nil
		}
	case UINT16:
		{
			err := tool.SetUint(d.Value, controlService.ValueType_VALUE_TYPE_UINT16)
			if err != nil {
				return nil, err
			}
			return &value, nil
		}
	case UINT32:
		{
			err := tool.SetUint(d.Value, controlService.ValueType_VALUE_TYPE_UINT32)
			if err != nil {
				return nil, err
			}
			return &value, nil
		}
	case UINT64:
		{
			err := tool.SetUint(d.Value, controlService.ValueType_VALUE_TYPE_UINT64)
			if err != nil {
				return nil, err
			}
			return &value, nil
		}
	case INT8:
		{
			err := tool.SetInt(d.Value, controlService.ValueType_VALUE_TYPE_INT8)
			if err != nil {
				return nil, err
			}
			return &value, nil
		}
	case INT16:
		{
			err := tool.SetInt(d.Value, controlService.ValueType_VALUE_TYPE_INT16)
			if err != nil {
				return nil, err
			}
			return &value, nil
		}
	case INT32:
		{
			err := tool.SetInt(d.Value, controlService.ValueType_VALUE_TYPE_INT32)
			if err != nil {
				return nil, err
			}
			return &value, nil
		}
	case INT64:
		{
			err := tool.SetInt(d.Value, controlService.ValueType_VALUE_TYPE_INT64)
			if err != nil {
				return nil, err
			}
			return &value, nil
		}
	case FLOAT32:
		{
			err := tool.SetFloat32(d.Value)
			if err != nil {
				return nil, err
			}
			return &value, nil
		}
	case FLOAT64:
		{
			err := tool.SetFloat64(d.Value)
			if err != nil {
				return nil, err
			}
			return &value, nil
		}
	case STRING:
		{
			err := tool.SetString(d.Value)
			if err != nil {
				return nil, err
			}
			return &value, nil
		}
	default:
		return nil, errors.New("unknown value	type")
	}
}
