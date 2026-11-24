package core

const DomainsWhoisXMLAPI = "https://domains-subdomains-discovery.whoisxmlapi.com/api/v1"

type Domains struct {
	Include []string `json:"include"`
}

type DomainData struct {
	APIKey  string  `json:"apiKey"`
	Domains Domains `json:"domains"`
}

func GetDomains(domain string, apiKey string) ([]string, error) {
	object, err := requestAPI(DomainsWhoisXMLAPI, "Domains & Subdomains Discovery API", apiKey, DomainData{
		APIKey:  apiKey,
		Domains: Domains{Include: []string{domain}},
	})
	if err != nil {
		return nil, err
	}
	return domainList(object, "Domains & Subdomains Discovery API")
}
