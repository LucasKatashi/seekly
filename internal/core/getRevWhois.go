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

const RevWhoisXMLAPI = "https://reverse-whois.whoisxmlapi.com/api/v2"

type BasicSearchTerms struct {
	Include []string `json:"include"`
}

type RevWhoisData struct {
	ApiKey           string           `json:"apiKey"`
	SearchType       string           `json:"searchType"`
	Mode             string           `json:"mode"`
	Punycode         bool             `json:"punycode"`
	BasicSearchTerms BasicSearchTerms `json:"basicSearchTerms"`
}

func GetRevWhois(queries []string, apiKey string) []string {
	allDomainsMap := make(map[string]struct{})

	for _, query := range queries {
		jsonParse := RevWhoisData{
			ApiKey:     apiKey,
			SearchType: "current",
			Mode:       "purchase",
			Punycode:   true,
			BasicSearchTerms: BasicSearchTerms{
				Include: []string{query},
			},
		}

		red := color.New(color.FgRed)

		jsonData, err := json.Marshal(jsonParse)
		if err != nil {
			fmt.Printf("[%s] Failed to marshal JSON object: %v\n", red.Sprint("ERR"), err)
			os.Exit(1)
		}

		resp, err := http.Post(RevWhoisXMLAPI, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("[%s] Failed to perform HTTP POST to RevWhoisXMLAPI (%s): %v\n", red.Sprint("ERR"), RevWhoisXMLAPI, err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("[%s] Failed to read response body from RevWhoisXMLAPI (%s): %v\n", red.Sprint("ERR"), RevWhoisXMLAPI, err)
			os.Exit(1)
		}

		domains := ExtractJsonValues(body)

		for _, d := range domains {
			allDomainsMap[d] = struct{}{}
		}
	}

	var allDomains []string
	for d := range allDomainsMap {
		allDomains = append(allDomains, d)
	}

	return allDomains
}
