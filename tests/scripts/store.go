package scripts

import "github.com/e154/smart-home/system/scripts"

var store interface{}

func storeRegisterCallback(scriptService *scripts.ScriptService) {
	scriptService.PushFunctions("store", func(value interface{}) {
		store = value
	})
}

type MyStruct struct {
	Bool    bool
	Int     int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	UInt    uint
	UInt8   uint8
	UInt16  uint16
	UInt32  uint32
	UInt64  uint64
	String  string
	Bytes   []byte
	Float32 float32
	Float64 float64
	Empty   *MyStruct
	Nested  *MyStruct
	Slice   []int
	private int
}

func (m *MyStruct) Multiply(x int) int {
	return m.Int * x
}

func (m *MyStruct) privateMethod() int {
	return 1
}
