package configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const separatorPattern string = "_"

// ExpandMap expands the map: ENV_PROPERTY_VALUE => {"ENV" : { "PROPERTY" : "VALUE"}}
func ExpandMap(mappings map[string]string) map[string]interface{} {
	expandedMap := make(map[string]interface{})
	for key, value := range mappings {
		expand(expandedMap, key, value)
	}
	return expandedMap
}

func after(value string, a string) string {
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:]
}

func expand(mapping map[string]interface{}, key string, value interface{}) {
	if strings.Contains(key, separatorPattern) {
		currentKey := strings.Split(key, separatorPattern)[0]
		if innerMap, ok := (mapping[currentKey]).(map[string]interface{}); ok {
			expand(innerMap, substring(key, separatorPattern), value)
		} else {
			innerMap := make(map[string]interface{})
			mapping[currentKey] = innerMap
			expand(innerMap, substring(key, separatorPattern), value)
		}
	} else {
		mapping[key] = value
	}
}

func substring(value string, delimeter string) string {
	index := strings.Index(value, delimeter)
	return value[index+1:]
}

func getEnvironment() map[string]string {
	settings := make(map[string]string)
	entries := os.Environ()
	for _, entry := range entries {
		split := strings.Split(entry, "=")
		settings[strings.ToLower(split[0])] = split[1]
	}
	return settings
}

// Bind the mapping to the configInterface
func Bind(mapping map[string]string, configInterface interface{}) (err error) {
	expandedMapping := ExpandMap(mapping)
	data, err := json.Marshal(expandedMapping)
	json.Unmarshal(data, configInterface)
	return
}

// BindEnvironment the configInterface to the environment variables
func BindEnvironment(configInterface interface{}) (err error) {
	err = Bind(getEnvironment(), configInterface)
	return
}

// InitializeConfig initialize config using the config tag
func InitializeConfig(mapping map[string]string, config interface{}) error {
	expandedMapping := ExpandMap(mapping)
	err := bindConfig(expandedMapping, config)
	if err != nil {
		return err
	}
	return nil
}

// BindConfig with data tag
func bindConfig(data map[string]interface{}, a interface{}) error {
	elementValue := reflect.ValueOf(a).Elem()
	elementType := reflect.TypeOf(a).Elem()

	for j := 0; j < elementValue.NumField(); j++ {
		field := elementValue.Field(j)
		fieldType := elementType.Field(j)
		configTag := elementValue.Type().Field(j).Tag.Get("config")
		if configTag == "" {
			return fmt.Errorf("no config tag")
		}
		value, ok := data[configTag]
		if !ok {
			return fmt.Errorf("no value found for config tag %s", configTag)
		}
		switch fieldType.Type.Kind() {
		case reflect.Bool:
			value, err := strconv.ParseBool(fmt.Sprint(value))
			if err != nil {
				return err
			}
			field.SetBool(value)
		case reflect.Struct:
			entry, ok := value.(map[string]interface{})
			if !ok {
				return fmt.Errorf("value for kind struct must be map[string]interface{}")
			}

			ptr := reflect.PtrTo(elementValue.Type().Field(j).Type)
			structure := reflect.New(ptr.Elem())
			field.Set(structure.Elem())
			err := bindConfig(entry, structure.Interface())
			if err != nil {
				return err
			}
			field.Set(structure.Elem())
		default:
			field.SetString(fmt.Sprint(value))
		}
	}
	return nil
}
