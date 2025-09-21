package envparse

import (
	"encoding/json"
	"fmt"
	"strings"
)

var parsers = make(map[string]func(string) (interface{}, error))

func RegisterParser(name string, parser func(string) (interface{}, error)) {
	parsers[name] = parser
}

func defaultSliceParser(value string) (interface{}, error) {
	return strings.Split(value, ","), nil
}

func defaultStructParser(value string) (interface{}, error) {
	var structValue map[string]interface{}
	err := json.Unmarshal([]byte(value), &structValue)
	if err != nil {
		return nil, err
	}

	return structValue, nil
}

func getParser(name string) (func(string) (interface{}, error), error) {
	found, ok := parsers[name]
	if !ok {
		return func(string) (interface{}, error) { return nil, nil }, fmt.Errorf("parser not found: %s", name)
	}

	return found, nil
}
