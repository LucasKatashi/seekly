package core

import (
	"encoding/json"
	"regexp"
	"strings"
)

var relevantContactFields = map[string]struct{}{
	"registrant":        {},
	"registrantContact": {},
}

var (
	ownerCRegex = regexp.MustCompile(`owner-c:\s*(\S+)`)
	nicHdlRegex = regexp.MustCompile(`nic-hdl-br:\s*(\S+)`)
	emailRegex  = regexp.MustCompile(`e-mail:\s*(\S+)`)
)

func extractOwnerEmailFromRawText(rawText string) string {
	ownerCMatch := ownerCRegex.FindStringSubmatch(rawText)
	if len(ownerCMatch) < 2 {
		return ""
	}
	ownerHandle := ownerCMatch[1]

	blocks := strings.Split(rawText, "nic-hdl-br:")
	for _, block := range blocks[1:] {
		nicMatch := nicHdlRegex.FindStringSubmatch("nic-hdl-br:" + block)
		if len(nicMatch) >= 2 && nicMatch[1] == ownerHandle {
			emailMatch := emailRegex.FindStringSubmatch(block)
			if len(emailMatch) >= 2 {
				return emailMatch[1]
			}
		}
	}
	return ""
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
						ownerEmail := extractOwnerEmailFromRawText(strValue)
						if ownerEmail != "" {
							uniqueStrings[ownerEmail] = struct{}{}
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
