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

const DomainsWhoisXMLAPI = "https://domains-subdomains-discovery.whoisxmlapi.com/api/v1"

type Domains struct {
	Include []string `json:"include"`
}

type DomainData struct {
	APIKey  string  `json:"apiKey"`
	Domains Domains `json:"domains"`
}

func GetDomains(domain string, apiKey string) []string {
	jsonObj := DomainData{
		APIKey: apiKey,
		Domains: Domains{
			Include: []string{domain},
		},
	}

	red := color.New(color.FgRed)

	jsonData, err := json.Marshal(jsonObj)
	if err != nil {
		fmt.Printf("[%s] Error parsing the JSON.", red.Sprint("ERR"))
		os.Exit(1)
	}

	resp, err := http.Post(DomainsWhoisXMLAPI, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("[%s] Failed to perform HTTP POST to DomainsWhoisXMLAPI (%s): %v\n", red.Sprint("ERR"), DomainsWhoisXMLAPI, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[%s] Failed to read response body from DomainsWhoisXMLAPI (%s): %v\n", red.Sprint("ERR"), DomainsWhoisXMLAPI, err)
		os.Exit(1)
	}

	values := ExtractJsonValues(body)

	return values
}
