package core

import "sort"

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

func GetRevWhois(queries []string, apiKey string) ([]string, error) {
	allDomainsMap := make(map[string]struct{})
	for _, query := range queries {
		object, err := requestAPI(RevWhoisXMLAPI, "Reverse WHOIS API", apiKey, RevWhoisData{
			ApiKey: apiKey, SearchType: "current", Mode: "purchase", Punycode: true,
			BasicSearchTerms: BasicSearchTerms{Include: []string{query}},
		})
		if err != nil {
			return nil, err
		}
		domains, err := domainList(object, "Reverse WHOIS API")
		if err != nil {
			return nil, err
		}
		for _, domain := range domains {
			allDomainsMap[domain] = struct{}{}
		}
	}
	allDomains := make([]string, 0, len(allDomainsMap))
	for domain := range allDomainsMap {
		allDomains = append(allDomains, domain)
	}
	sort.Strings(allDomains)
	return allDomains, nil
}
