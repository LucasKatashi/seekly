package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/fatih/color"
)

const WhoisXMLAPI = "https://www.whoisxmlapi.com/whoisserver/WhoisService"

type WhoisData struct {
	APIKey         string `json:"apiKey"`
	DomainName     string `json:"domainName"`
	OutputFormat   string `json:"outputFormat"`
	IgnoreRawTexts int    `json:"ignoreRawTexts"`
}

func GetWhois(domain string, apiKey string) []string {
	jsonObj := WhoisData{
		APIKey:         apiKey,
		DomainName:     domain,
		OutputFormat:   "JSON",
		IgnoreRawTexts: 1,
	}

	red := color.New(color.FgRed)

	jsonData, err := json.Marshal(jsonObj)
	if err != nil {
		fmt.Printf("[%s] Failed to marshal JSON object: %v\n", red.Sprint("ERR"), err)
		os.Exit(1)
	}

	resp, err := http.Post(WhoisXMLAPI, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("[%s] Error querying the Whois API %s: %v.\n", red.Sprint("ERR"), WhoisXMLAPI, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[%s] Failed to read response from %s: %v\n", red.Sprint("ERR"), WhoisXMLAPI, err)
		os.Exit(1)
	}

	values := ExtractJsonValues(body)

	return values
}
