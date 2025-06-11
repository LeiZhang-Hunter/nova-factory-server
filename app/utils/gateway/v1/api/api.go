package api

type DataValueType string

var Bit DataValueType = "bit"
var UInt8 DataValueType = "uint8"
var UInt16 DataValueType = "uint16"
var Int16 DataValueType = "int16"
var UInt32 DataValueType = "uint32"
var Int32 DataValueType = "int32"
var Float32 DataValueType = "float32"
var Float64 DataValueType = "float64"

type ByteOrder string

var BigEndian ByteOrder = "ABCD"
var LittleEndian ByteOrder = "CDAB"
