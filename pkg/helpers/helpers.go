package helpers

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// CompareJSONWithMap compares map and string as json
func CompareJSONWithMap(jsonStr string, m map[string]interface{}) (bool, error) {
	// Unmarshal the JSON string into a map
	var jsonMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonMap)
	if err != nil {
		return false, fmt.Errorf("error unmarshalling JSON string: %v", err)
	}

	// Marshal both maps into JSON to normalize them
	jsonMapBytes, err1 := json.Marshal(jsonMap)
	if err1 != nil {
		return false, fmt.Errorf("error marshalling jsonMap to JSON: %v", err1)
	}

	mBytes, err2 := json.Marshal(m)
	if err2 != nil {
		return false, fmt.Errorf("error marshalling map to JSON: %v", err2)
	}

	// Compare the JSON byte slices
	areEqual := reflect.DeepEqual(jsonMapBytes, mBytes)
	return areEqual, nil
}

// CompareMapsAsJSON compares two maps as JSON
func CompareMapsAsJSON(map1, map2 map[string]interface{}) (bool, error) {
	// Marshal the first map into JSON
	json1, err1 := json.Marshal(map1)
	if err1 != nil {
		return false, fmt.Errorf("error marshalling map1 to JSON: %v", err1)
	}

	// Marshal the second map into JSON
	json2, err2 := json.Marshal(map2)
	if err2 != nil {
		return false, fmt.Errorf("error marshalling map2 to JSON: %v", err2)
	}

	// Compare the JSON byte slices
	areEqual := reflect.DeepEqual(json1, json2)
	return areEqual, nil
}

func CompareJSONStrings(jsonStr1, jsonStr2 string) (bool, error) {
	// Unmarshal the first JSON string into a map
	var map1 map[string]interface{}
	err1 := json.Unmarshal([]byte(jsonStr1), &map1)
	if err1 != nil {
		return false, fmt.Errorf("error unmarshalling jsonStr1: %v", err1)
	}

	// Unmarshal the second JSON string into a map
	var map2 map[string]interface{}
	err2 := json.Unmarshal([]byte(jsonStr2), &map2)
	if err2 != nil {
		return false, fmt.Errorf("error unmarshalling jsonStr2: %v", err2)
	}

	// Compare the two maps
	areEqual := reflect.DeepEqual(map1, map2)
	return areEqual, nil
}
