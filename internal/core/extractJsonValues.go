package core

import (
	"encoding/json"
)

func ExtractJsonValues(jsonData []byte) []string {
	var decodedData interface{}
	if err := json.Unmarshal(jsonData, &decodedData); err != nil {
		return nil
	}

	uniqueStrings := make(map[string]struct{})

	var traverseData func(data interface{})

	traverseData = func(data interface{}) {
		switch v := data.(type) {

		case map[string]interface{}:
			for key, value := range v {
				if key == "name" || key == "email" || key == "contactEmail" || key == "domainsList" {
					if key == "domainsList" {
						if arr, ok := value.([]interface{}); ok {
							for _, item := range arr {
								if s, ok := item.(string); ok {
									uniqueStrings[s] = struct{}{}
								}
							}
						}
					} else {
						if strValue, ok := value.(string); ok {
							uniqueStrings[strValue] = struct{}{}
						}
					}
				}
				traverseData(value)
			}
		case []interface{}:
			for _, item := range v {
				traverseData(item)
			}

		}

	}
	traverseData(decodedData)

	var result []string
	for str := range uniqueStrings {
		result = append(result, str)
	}

	return result
}
