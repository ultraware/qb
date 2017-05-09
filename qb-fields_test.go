package qb

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testValue struct {
	Data     interface{}
	Zero     interface{}
	Function func(Field) DataField
}

var testValues = []testValue{
	{
		Data:     7,
		Zero:     0,
		Function: func(f Field) DataField { return NewIntField(f) },
	},
	{
		Data:     int32(7),
		Zero:     int32(0),
		Function: func(f Field) DataField { return NewInt32Field(f) },
	},
	{
		Data:     int64(7),
		Zero:     int64(0),
		Function: func(f Field) DataField { return NewInt64Field(f) },
	},
	{
		Data:     `Hello, world!`,
		Zero:     ``,
		Function: func(f Field) DataField { return NewStringField(f) },
	},
	{
		Data:     []byte(`Hello, world!`),
		Zero:     []byte(nil),
		Function: func(f Field) DataField { return NewBytesField(f) },
	},
	{
		Data:     float64(3.14159265359),
		Zero:     float64(0),
		Function: func(f Field) DataField { return NewFloat64Field(f) },
	},
	{
		Data:     float32(3.14159265359),
		Zero:     float32(0),
		Function: func(f Field) DataField { return NewFloat32Field(f) },
	},
	{
		Data:     true,
		Zero:     false,
		Function: func(f Field) DataField { return NewBoolField(f) },
	},
	{
		Data:     time.Now(),
		Zero:     time.Time{},
		Function: func(f Field) DataField { return NewTimeField(f) },
	},

	/////// Null fields ///////
	{
		Data:     getPointer(7),
		Zero:     nil,
		Function: func(f Field) DataField { return NewNullIntField(f) },
	},
	{
		Data:     getPointer(int32(7)),
		Zero:     nil,
		Function: func(f Field) DataField { return NewNullInt32Field(f) },
	},
	{
		Data:     getPointer(int64(7)),
		Zero:     nil,
		Function: func(f Field) DataField { return NewNullInt64Field(f) },
	},
	{
		Data:     getPointer(`Hello, world!`),
		Zero:     nil,
		Function: func(f Field) DataField { return NewNullStringField(f) },
	},
	{
		Data:     getPointer([]byte(`Hello, world!`)),
		Zero:     nil,
		Function: func(f Field) DataField { return NewNullBytesField(f) },
	},
	{
		Data:     getPointer(float64(3.14159265359)),
		Zero:     nil,
		Function: func(f Field) DataField { return NewNullFloat64Field(f) },
	},
	{
		Data:     getPointer(float32(3.14159265359)),
		Zero:     nil,
		Function: func(f Field) DataField { return NewNullFloat32Field(f) },
	},
	{
		Data:     getPointer(true),
		Zero:     nil,
		Function: func(f Field) DataField { return NewNullBoolField(f) },
	},
	{
		Data:     getPointer(time.Now()),
		Zero:     nil,
		Function: func(f Field) DataField { return NewNullTimeField(f) },
	},
}

func TestAll(t *testing.T) {
	assert := assert.New(t)

	testField := TableField{Name: `abc`}

	for i := 0; i < len(testValues); i++ {
		var (
			testValue = testValues[i]

			data  = testValue.Data
			zero  = testValue.Zero
			field = testValue.Function(testField)

			input = reflect.ValueOf(data)
			val   = reflect.ValueOf(field)
		)

		t.Logf(`Testing: %T`, field)

		assert.Equal(false, field.hasChanged())
		assert.Equal(false, field.isSet())

		assert.Equal(testField.Name, field.getField().(TableField).Name)

		val.MethodByName(`Set`).Call([]reflect.Value{input})
		assert.Equal(true, field.hasChanged())
		assert.Equal(true, field.isSet())

		field.Reset()
		assert.Equal(false, field.hasChanged())
		assert.Equal(false, field.isSet())
		out, err := field.Value()
		assert.NoError(err)
		if zero == nil {
			assert.Nil(out)
		} else {
			assert.Equal(zero, out)
		}

		// Reset field
		field = testValue.Function(testField)
		val = reflect.ValueOf(field)

		err = field.Scan(getDriverValue(data))
		assert.NoError(err)
		assert.Equal(true, field.isSet())
		assert.Equal(false, field.hasChanged())

		output := val.MethodByName(`Get`).Call(nil)
		assert.Equal(reflect.ValueOf(data).Interface(), output[0].Interface())

		if zero == nil {
			field.Scan(zero)

			output := val.MethodByName(`Get`).Call(nil)
			assert.Nil(output[0].Interface())
		}
	}
}

func getDriverValue(in interface{}) interface{} {
	switch v := in.(type) {
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case float32:
		return float64(v)

	case *int:
		return int64(*v)
	case *int32:
		return int64(*v)
	case *int64:
		return *v
	case *float32:
		return float64(*v)
	case *float64:
		return *v
	case *string:
		return *v
	case *[]byte:
		return *v
	case *bool:
		return *v
	case *time.Time:
		return *v

	default:
		return v
	}
}

func getPointer(in interface{}) interface{} {
	switch v := in.(type) {
	case int:
		return &v
	case int64:
		return &v
	case int32:
		return &v
	case string:
		return &v
	case float32:
		return &v
	case float64:
		return &v
	case bool:
		return &v
	case []byte:
		return &v
	case time.Time:
		return &v
	default:
		fmt.Printf(`Unknown type: %T`, v)
		return &in
	}
}
