package envparse

import (
	"reflect"
	"strings"
	"testing"
)

type Primatives struct {
	TestBool    bool                   `envparse:"TEST_BOOL"`
	TestInt     int                    `envparse:"TEST_INT"`
	TestInt8    int8                   `envparse:"TEST_INT"`
	TestInt16   int16                  `envparse:"TEST_INT"`
	TestInt32   int32                  `envparse:"TEST_INT"`
	TestInt64   int64                  `envparse:"TEST_INT"`
	TestUint    uint                   `envparse:"TEST_UINT"`
	TestUint8   uint8                  `envparse:"TEST_UINT"`
	TestUint16  uint16                 `envparse:"TEST_UINT"`
	TestUint32  uint32                 `envparse:"TEST_UINT"`
	TestUint64  uint64                 `envparse:"TEST_UINT"`
	TestFloat32 float32                `envparse:"TEST_FLOAT"`
	TestFloat64 float64                `envparse:"TEST_FLOAT"`
	TestString  string                 `envparse:"TEST_STRING"`
	TestSlice   []string               `envparse:"name=TEST_SLICE,parser=default"`
	TestMap     map[string]interface{} `envparse:"name=TEST_MAP,parser=default"`
	DoNothing   string
}

type MissingEnv struct {
	TestInt int `envparse:"MISSING"`
}

type BadType struct {
	TestBadType interface{} `envparse:"TEST_MAP"`
}

type CustomParserType struct {
	CustomMap   map[string]string `envparse:"name=TEST_CUSTOM_MAP,parser=mapEquals"`
	CustomSlice []int             `envparse:"name=TEST_SLICE,parser=intSlice"`
}

type BadCustomMapParserType struct {
	CustomMap map[string]string `envparse:"name=TEST_CUSTOM_MAP,parser=noParser"`
}

type BadCustomSliceParserType struct {
	CustomSlice []int `envparse:"name=TEST_SLICE,parser=noParser"`
}

func TestParseTagValues(t *testing.T) {
	t.Run("only env name", func(t *testing.T) {
		values := parseTagValues("TEST_ENV")
		if values.envName != "TEST_ENV" {
			t.Errorf("Expected envName to be TEST_ENV, got %s", values.envName)
		}
	})

	t.Run("env name with parser", func(t *testing.T) {
		values := parseTagValues("name=TEST_ENV,parser=string")
		if values.parserName != "string" {
			t.Errorf("Expected parserName to be string, got %s", values.parserName)
		}
	})
}

func TestParse(t *testing.T) {
	t.Run("primatives", func(t *testing.T) {
		var p Primatives
		err := Parse(&p)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if p.TestBool != true {
			t.Errorf("Expected testBool to be true, got %v", p.TestBool)
		}

		if p.TestInt != -1 || p.TestInt8 != -1 || p.TestInt16 != -1 || p.TestInt32 != -1 || p.TestInt64 != -1 {
			t.Errorf("Expected all ints to be -1, got %v, %v, %v, %v, %v", p.TestInt, p.TestInt8, p.TestInt16, p.TestInt32, p.TestInt64)
		}

		if p.TestUint != 1 || p.TestUint8 != 1 || p.TestUint16 != 1 || p.TestUint32 != 1 || p.TestUint64 != 1 {
			t.Errorf("Expected all uints to be 1, got %v, %v, %v, %v, %v", p.TestUint, p.TestUint8, p.TestUint16, p.TestUint32, p.TestUint64)
		}

		if p.TestFloat32 != 0.1 || p.TestFloat64 != 0.1 {
			t.Errorf("Expected floats to be 0.1, got %v, %v", p.TestFloat32, p.TestFloat64)
		}

		if p.TestString != "test" {
			t.Errorf("Expected testString to be test, got %v", p.TestString)
		}

		if len(p.TestSlice) != 2 {
			t.Errorf("Expected TestSlice to have length of 2, got %v", len(p.TestSlice))
		}

		if value, ok := p.TestMap["test"]; !ok || value != "test" {
			t.Errorf("Expected TestMap[\"test\"] to equal \"test\" got %v", value)
		}
	})

	t.Run("missing env var", func(t *testing.T) {
		var p MissingEnv
		err := Parse(&p)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if err.Error() != "unable to find env var MISSING" {
			t.Errorf("Expected error to be 'unable to find env var MISSING', got %v", err.Error())
		}
	})

	t.Run("test object with struct", func(t *testing.T) {
		var p BadType
		err := Parse(&p)
		if err == nil {
			t.Error("Expected error, got nil")
		}

		if err.Error() != "no registed set method for type" {
			t.Errorf("Got unexpected error message %v", err)
		}
	})
}

func TestCustomParsers(t *testing.T) {
	t.Run("Test good parsers", func(t *testing.T) {
		customMapParser := func(value string) (interface{}, error) {
			result := make(map[string]string)

			splitValue := strings.Split(value, "=")

			result[splitValue[0]] = splitValue[1]
			return result, nil
		}

		customSliceParser := func(value string) (interface{}, error) {

			splitValue := strings.Split(value, ",")

			result := make([]int, 0, len(splitValue))

			for i := range splitValue {
				result = append(result, i)
			}

			return result, nil
		}

		RegisterParser("mapEquals", customMapParser)
		RegisterParser("intSlice", customSliceParser)

		var p CustomParserType
		err := Parse(&p)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if value, ok := p.CustomMap["test"]; !ok || value != "test" {
			t.Errorf("Expected p.CustomMap[\"test\"] to equal 'test' got %v", value)
		}

		if len(p.CustomSlice) != 2 {
			t.Errorf("Expected p.CustomSlice to have length of 2 got %v", len(p.CustomSlice))
		}
	})

	t.Run("Test bad map parser", func(t *testing.T) {
		var p BadCustomMapParserType
		err := Parse(&p)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("Test bad slice parser", func(t *testing.T) {
		var p BadCustomSliceParserType
		err := Parse(&p)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

}

func TestSetField(t *testing.T) {
	type Test struct {
		Field string
	}

	test := Test{}

	field := reflect.ValueOf(&test).Elem().Type().Field(0)

	err := setField(&test, field, "test", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if test.Field != "test" {
		t.Errorf("Expected Field to be test, got %v", test.Field)
	}
}
