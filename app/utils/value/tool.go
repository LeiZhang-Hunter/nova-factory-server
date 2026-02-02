package value

import (
	v1 "github.com/novawatcher-io/nova-factory-payload/control/v1"
	"github.com/pkg/errors"
)

type Tool struct {
	ptr *v1.Value
}

func NewValueTool(ptr *v1.Value) *Tool {
	return &Tool{
		ptr: ptr,
	}
}

func (v *Tool) SetBit(value Value) error {
	v.ptr.Type = v1.ValueType_VALUE_TYPE_BIT
	if v.ptr.GetValue() == nil {
		v.ptr.Value = &v1.Value_ByteValue{}
	}
	x, ok := v.ptr.GetValue().(*v1.Value_ByteValue)
	if !ok {
		return errors.New("value is not string")
	}
	x.ByteValue = value.ByteValue
	return nil
}

func (v *Tool) SetUint(value Value, vType v1.ValueType) error {
	v.ptr.Type = vType
	if v.ptr.GetValue() == nil {
		v.ptr.Value = &v1.Value_UintValue{}
	}
	x, ok := v.ptr.GetValue().(*v1.Value_UintValue)
	if !ok {
		return errors.New("value is not string")
	}
	x.UintValue = value.UintValue
	return nil
}

func (v *Tool) SetInt(value Value, vType v1.ValueType) error {
	v.ptr.Type = vType
	if v.ptr.GetValue() == nil {
		v.ptr.Value = &v1.Value_IntValue{}
	}
	x, ok := v.ptr.GetValue().(*v1.Value_IntValue)
	if !ok {
		return errors.New("value is not string")
	}
	x.IntValue = value.IntValue
	return nil
}

func (v *Tool) SetFloat32(value Value) error {
	v.ptr.Type = v1.ValueType_VALUE_TYPE_FLOAT32
	if v.ptr.GetValue() == nil {
		v.ptr.Value = &v1.Value_FloatValue{}
	}
	x, ok := v.ptr.GetValue().(*v1.Value_FloatValue)
	if !ok {
		return errors.New("value is not string")
	}
	x.FloatValue = value.FloatValue
	return nil
}

func (v *Tool) SetFloat64(value Value) error {
	v.ptr.Type = v1.ValueType_VALUE_TYPE_FLOAT64
	if v.ptr.GetValue() == nil {
		v.ptr.Value = &v1.Value_FloatValue{}
	}
	x, ok := v.ptr.GetValue().(*v1.Value_FloatValue)
	if !ok {
		return errors.New("value is not string")
	}
	x.FloatValue = value.FloatValue
	return nil
}

func (v *Tool) SetString(value Value) error {
	v.ptr.Type = v1.ValueType_VALUE_TYPE_STRING
	if v.ptr.GetValue() == nil {
		v.ptr.Value = &v1.Value_StringValue{}
	}
	x, ok := v.ptr.GetValue().(*v1.Value_StringValue)
	if !ok {
		return errors.New("value is not string")
	}
	x.StringValue = value.StringValue
	return nil
}
