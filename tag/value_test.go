package tag

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValue_Set(t *testing.T) {
	type TestStruct struct {
		IntField     int     `test:"intField"`
		Int8Field    int8    `test:"int8Field"`
		Int16Field   int16   `test:"int16Field"`
		Int32Field   int32   `test:"int32Field"`
		Int64Field   int64   `test:"int64Field"`
		UintField    uint    `test:"uintField"`
		Uint8Field   uint8   `test:"uint8Field"`
		Uint16Field  uint16  `test:"uint16Field"`
		Uint32Field  uint32  `test:"uint32Field"`
		Uint64Field  uint64  `test:"uint64Field"`
		Float32Field float32 `test:"float32Field"`
		Float64Field float64 `test:"float64Field"`
		StringField  string  `test:"stringField"`
		BoolField    bool    `test:"boolField"`
		PtrField     *string `test:"ptrField"`
	}

	var testStruct TestStruct
	parser := NewParser(parseFunc)
	values := parser.Parse("test", &testStruct)
	values.Get("intField").Set(1)
	values.Get("int8Field").Set(int8(1))
	values.Get("int16Field").Set(int16(1))
	values.Get("int32Field").Set(int32(1))
	values.Get("int64Field").Set(int64(1))
	values.Get("uintField").Set(uint(1))
	values.Get("uint8Field").Set(uint8(1))
	values.Get("uint16Field").Set(uint16(1))
	values.Get("uint32Field").Set(uint32(1))
	values.Get("uint64Field").Set(uint64(1))
	values.Get("float32Field").Set(float32(1.0))
	values.Get("float64Field").Set(float64(1.0))
	values.Get("stringField").Set("1")
	values.Get("boolField").Set(true)
	assert.Equal(t, 1, testStruct.IntField)
	assert.Equal(t, int8(1), testStruct.Int8Field)
	assert.Equal(t, int16(1), testStruct.Int16Field)
	assert.Equal(t, int32(1), testStruct.Int32Field)
	assert.Equal(t, int64(1), testStruct.Int64Field)
	assert.Equal(t, uint(1), testStruct.UintField)
	assert.Equal(t, uint8(1), testStruct.Uint8Field)
	assert.Equal(t, uint16(1), testStruct.Uint16Field)
	assert.Equal(t, uint32(1), testStruct.Uint32Field)
	assert.Equal(t, uint64(1), testStruct.Uint64Field)
	assert.Equal(t, float32(1.0), testStruct.Float32Field)
	assert.Equal(t, float64(1.0), testStruct.Float64Field)
	assert.Equal(t, "1", testStruct.StringField)
	assert.Equal(t, true, testStruct.BoolField)

	// int uint conversion
	values.Get("intField").Set(uint(2))
	values.Get("int8Field").Set(uint8(2))
	values.Get("int16Field").Set(uint16(2))
	values.Get("int32Field").Set(uint32(2))
	values.Get("int64Field").Set(uint64(2))
	values.Get("uintField").Set(int(2))
	values.Get("uint8Field").Set(int8(2))
	values.Get("uint16Field").Set(int16(2))
	values.Get("uint32Field").Set(int32(2))
	values.Get("uint64Field").Set(int64(2))
	assert.Equal(t, int(2), testStruct.IntField)
	assert.Equal(t, int8(2), testStruct.Int8Field)
	assert.Equal(t, int16(2), testStruct.Int16Field)
	assert.Equal(t, int32(2), testStruct.Int32Field)
	assert.Equal(t, int64(2), testStruct.Int64Field)
	assert.Equal(t, uint(2), testStruct.UintField)
	assert.Equal(t, uint8(2), testStruct.Uint8Field)
	assert.Equal(t, uint16(2), testStruct.Uint16Field)
	assert.Equal(t, uint32(2), testStruct.Uint32Field)
	assert.Equal(t, uint64(2), testStruct.Uint64Field)

	// string conversion
	values.Get("stringField").Set([]byte("2"))
	assert.Equal(t, "2", testStruct.StringField)
	values.Get("stringField").Set([]rune("3"))
	assert.Equal(t, "3", testStruct.StringField)

	// fallback int/uint
	type CustomInt int
	values.Get("intField").Set(CustomInt(3))
	assert.Equal(t, 3, testStruct.IntField)
	values.Get("uintField").Set(CustomInt(3))
	assert.Equal(t, uint(3), testStruct.UintField)

	// fallback float
	type CustomFloat float32
	values.Get("float32Field").Set(CustomFloat(3.0))
	assert.Equal(t, float32(3.0), testStruct.Float32Field)

	// fallback string
	type CustomString string
	values.Get("stringField").Set(CustomString("4"))
	assert.Equal(t, "4", testStruct.StringField)

	// fallback bool
	type CustomBool bool
	values.Get("boolField").Set(CustomBool(false))
	assert.Equal(t, false, testStruct.BoolField)

	// set nil
	var stringValue string = "haha"
	testStruct.PtrField = &stringValue
	values.Get("ptrField").Set(nil)
	assert.Nil(t, testStruct.PtrField)
}
