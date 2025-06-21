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

func CheckByteOrder(order string) bool {
	if order != string(BigEndian) && order != string(LittleEndian) {
		return false
	}
	return true
}

var StorageType string
var ChangeStorage string = "change_storage"
var peerStorage string = "peer_storage"
var peer30sStorage string = "peer_30s_storage"
var peer1minStorage string = "peer_1min_storage"
var peer5minStorage string = "peer_5min_storage"
var peer10minStorage string = "peer_10min_storage"
var peer30minStorage string = "peer_30min_storage"
var peer1hStorage string = "peer_1h_storage"
var peer1dayStorage string = "peer_1day_storage"

func CheckStorageType(storageType string) bool {
	if storageType != peerStorage && storageType != peer30sStorage && storageType != peer1minStorage &&
		storageType != ChangeStorage && storageType != peer5minStorage && storageType != peer10minStorage &&
		storageType != peer30minStorage && storageType != peer1hStorage && StorageType != peer1dayStorage {
		return false
	}
	return true
}

var FunctionCode string

const (
	// Bit access
	FuncCodeReadDiscreteInputs = 2
	FuncCodeReadCoils          = 1
	FuncCodeWriteSingleCoil    = 5
	FuncCodeWriteMultipleCoils = 15

	// 16-bit access
	FuncCodeReadInputRegisters         = 4
	FuncCodeReadHoldingRegisters       = 3
	FuncCodeWriteSingleRegister        = 6
	FuncCodeWriteMultipleRegisters     = 16
	FuncCodeReadWriteMultipleRegisters = 23
	FuncCodeMaskWriteRegister          = 22
	FuncCodeReadFIFOQueue              = 24
)

func CheckFunctionCode(functionCode int) bool {
	switch functionCode {
	case FuncCodeReadDiscreteInputs:
		return true
	case FuncCodeReadCoils:
		return true
	case FuncCodeWriteSingleCoil:
		return true
	case FuncCodeWriteMultipleCoils:
		return true
	case FuncCodeReadInputRegisters:
		return true
	case FuncCodeReadHoldingRegisters:
		return true
	case FuncCodeWriteSingleRegister:
		return true
	case FuncCodeWriteMultipleRegisters:
		return true
	case FuncCodeReadWriteMultipleRegisters:
		return true
	case FuncCodeMaskWriteRegister:
		return true
	case FuncCodeReadFIFOQueue:
		return true
	}
	return false
}
