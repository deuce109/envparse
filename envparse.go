package envparse

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type tagValues struct {
	parserName string
	envName    string
}

func parseTagValues(tagString string) tagValues {
	splitValues := strings.Split(tagString, ",")

	values := tagValues{}

	for _, v := range splitValues {
		splitValue := strings.Split(v, "=")
		switch splitValue[0] {
		case "name":
			values.envName = splitValue[1]
		case "parser":
			values.parserName = splitValue[1]
		}
	}

	if values.envName == "" {
		values.envName = tagString
	}

	return values
}

func Parse[T interface{}](obj *T) error {
	t := reflect.TypeOf(*obj)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagValue := field.Tag.Get("envparse")

		if tagValue == "" {
			continue
		}

		values := parseTagValues(tagValue)

		if value := os.Getenv(values.envName); value != "" {
			err := setField(obj, field, value, values.parserName)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unable to find env var %s", values.envName)
		}
	}

	return nil
}
