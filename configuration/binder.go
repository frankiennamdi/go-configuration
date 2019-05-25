package configuration

import (
	"encoding/json"
	"os"
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
	return value[adjustedPos:len(value)]
}

func expand(mapping map[string]interface{}, key string, value interface{}) {
	if strings.Contains(key, separatorPattern) {
		currentKey := strings.Split(key, separatorPattern)[0]
		if innerMap, ok := (mapping[currentKey]).(map[string]interface{}); ok {
			expand(innerMap, after(key, separatorPattern), value)
		} else {
			innerMap := make(map[string]interface{})
			mapping[currentKey] = innerMap
			expand(innerMap, after(key, separatorPattern), value)
		}
	} else {
		mapping[key] = value
	}
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
