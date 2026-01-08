package core

import (
	"encoding/json"
	"regexp"
)

var emailRegex = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)

var relevantContactFields = map[string]struct{}{
	"registrant":             {},
	"registrantContact":      {},
	"administrativeContact":  {},
	"technicalContact":       {},
	"billingContact":         {},
	"zoneContact":            {},
}

func ExtractJsonValues(jsonData []byte) []string {
	var decodedData interface{}
	if err := json.Unmarshal(jsonData, &decodedData); err != nil {
		return nil
	}

	uniqueStrings := make(map[string]struct{})

	var extractFromContact func(contact interface{})
	extractFromContact = func(contact interface{}) {
		contactMap, ok := contact.(map[string]interface{})
		if !ok {
			return
		}

		if name, ok := contactMap["name"].(string); ok && name != "" {
			uniqueStrings[name] = struct{}{}
		}

		if org, ok := contactMap["organization"].(string); ok && org != "" {
			uniqueStrings[org] = struct{}{}
		}

		if email, ok := contactMap["email"].(string); ok && email != "" {
			uniqueStrings[email] = struct{}{}
		}
	}

	var traverseData func(data interface{})

	traverseData = func(data interface{}) {
		switch v := data.(type) {

		case map[string]interface{}:
			for key, value := range v {
				if key == "domainsList" {
					if arr, ok := value.([]interface{}); ok {
						for _, item := range arr {
							if s, ok := item.(string); ok {
								uniqueStrings[s] = struct{}{}
							}
						}
					}
				}

				if _, isRelevant := relevantContactFields[key]; isRelevant {
					extractFromContact(value)
				}
				if key == "rawText" {
					if strValue, ok := value.(string); ok {
						emails := emailRegex.FindAllString(strValue, -1)
						for _, email := range emails {
							uniqueStrings[email] = struct{}{}
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
