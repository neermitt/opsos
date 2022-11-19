package utils

import (
	"encoding/json"
	"sort"
)

// StringKeysFromMap returns a slice of sorted string keys from the provided map
func StringKeysFromMap(m map[string]any) []string {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// MapKeyExists checks if a key already defined in a map
func MapKeyExists(m map[string]any, key string) bool {
	_, ok := m[key]
	return ok
}

func ToMap(data any) (map[string]any, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	err = json.Unmarshal(raw, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func FromMap(data map[string]any, target any) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, target)
}
