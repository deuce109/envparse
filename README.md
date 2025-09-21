
## EnvParse

### Usage

A struct with all of the currently possible value types is at the end of this file

Before running the below code sample please ensure you have an env variable of `DEMO_RAN` set to `true`

```go

package main

import (
	"errors"
	"github.com/deuce109/envparse/v2"
)

type Demo struct {
	Ran bool `envparse:"DEMO_RAN"`
}

func main() {

	var d Demo
	err := Parse(&d)

	if err != nil {
		panic(err)
	}

	if !d.Ran {
		panic(errors.New("d.Ran should have been true"))
	}

	println(d.Ran)

}


```

#### Tag
The `envparse` tag is used to specify what environment variable you would like to read into your struct fields.

It also allows for using custom parsing for slices and maps this is done by modifying the tag to be in the format `envparse:name=<env_name>,parser=<parser_name>`.

#### Custom Parsing

In order to register a custom parser the parser function must be of type `func (string) (interface{}, error)`.

Then you call `RegisterParser` with it's name and the function

This *must* be done before the call to `Parse` relying on the parser function.

```go

customSliceParser := func(value string) (interface{}, error) {

    splitValue := strings.Split(value, ",")

    result := make([]int, 0, len(splitValue))

    for i := range splitValue {
        result = append(result, i)
    }

    return result, nil
}

RegisterParser("intSlice", customSliceParser)

```

#### How
This package reflects over the given object in order to determine which type to attempt to parse to.

#### File Structure
`envparse.go`: Contains the main Parse function and it's helpers
`setters.go`: Contains the set methods for each type and some helper functions
`parsers.go`: Contains the ability for parsing of maps and slices from streams


#### Example Struct

The following struct has all of the available types that `envparse` can currently parse

```go
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
```