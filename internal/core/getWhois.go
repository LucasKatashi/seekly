package core

import (
	"encoding/json"
	"fmt"
	"strings"
)

const WhoisXMLAPI = "https://www.whoisxmlapi.com/whoisserver/WhoisService"

type WhoisData struct {
	APIKey         string `json:"apiKey"`
	DomainName     string `json:"domainName"`
	OutputFormat   string `json:"outputFormat"`
	IgnoreRawTexts int    `json:"ignoreRawTexts"`
}

func GetWhois(domain string, apiKey string) ([]string, error) {
	object, err := requestAPI(WhoisXMLAPI, "WHOIS API", apiKey, WhoisData{
		APIKey: apiKey, DomainName: domain, OutputFormat: "JSON", IgnoreRawTexts: 0,
	})
	if err != nil {
		return nil, err
	}
	record := object["WhoisRecord"]
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(record, &fields); err != nil || fields == nil {
		return nil, fmt.Errorf("WHOIS API: response is missing a valid WhoisRecord; cannot determine registrant details")
	}
	var dataError string
	if raw := fields["dataError"]; len(raw) > 0 {
		if err := json.Unmarshal(raw, &dataError); err != nil {
			return nil, fmt.Errorf("WHOIS API: invalid dataError field")
		}
		if dataError != "" {
			return nil, fmt.Errorf("WHOIS API: no usable WHOIS data: %s", strings.ReplaceAll(dataError, apiKey, "[REDACTED]"))
		}
	}
	return ExtractJsonValues(record), nil
}
