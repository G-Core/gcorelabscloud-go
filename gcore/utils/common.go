package utils

import "fmt"

func MapInterfaceToMapString(mapInterface interface{}) (map[string]string, error) {
	mapString := make(map[string]string)

	switch mapInterface.(type) {
	default:
		return nil, fmt.Errorf("Unexpected type %T", mapInterface)
	case map[string]interface{}:
		for key, value := range mapInterface.(map[string]interface{}) {
			mapString[key] = fmt.Sprintf("%v", value)
		}
	case map[interface{}]interface{}:
		for key, value := range mapInterface.(map[interface{}]interface{}) {
			mapString[fmt.Sprintf("%v", key)] = fmt.Sprintf("%v", value)
		}
	}

	return mapString, nil
}
