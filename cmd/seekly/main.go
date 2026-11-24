package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/LucasKatashi/seekly/internal/core"
	"github.com/LucasKatashi/seekly/internal/ui"
)

func main() {
	domain := flag.String("domain", "", "")
	silent := flag.Bool("silent", false, "")
	wildcard := flag.Bool("wildcard", false, "")
	apiKey := flag.String("api-key", "", "")
	output := flag.String("output", "", "")

	fail := func(err error) {
		fmt.Fprintf(os.Stderr, "[ERR] %v\n", err)
		os.Exit(1)
	}
	flag.Usage = ui.CustomUsage
	flag.Parse()
	ui.PrintBanner(*silent)

	if *apiKey == "" {
		*apiKey = os.Getenv("WhoisXMLAPIKey")
		if *apiKey == "" {
			fail(fmt.Errorf("API key is missing: set WhoisXMLAPIKey or use --api-key"))
		}
	}

	*domain = strings.TrimSuffix(strings.TrimSpace(*domain), ".")
	parts := strings.SplitN(*domain, ".", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" || strings.ContainsAny(*domain, "/: \t\n") {
		fail(fmt.Errorf("--domain must be a domain name such as example.com (not a URL)"))
	}
	if err := core.CheckAPIKey(*apiKey); err != nil {
		fail(err)
	}
	values, err := core.GetWhois(*domain, *apiKey)
	if err != nil {
		fail(fmt.Errorf("no domains returned: WHOIS lookup failed: %w", err))
	}
	var outputDomains []string
	if len(values) == 0 {
		fmt.Fprintln(os.Stderr, "[WARN] WHOIS returned no usable registrant name, organization, or email. These fields may be redacted, private, or absent; Reverse WHOIS was skipped.")
	} else {
		outputDomains, err = core.GetRevWhois(values, *apiKey)
		if err != nil {
			fail(fmt.Errorf("enumeration failed; results may be incomplete: %w", err))
		}
		if len(outputDomains) == 0 {
			fmt.Fprintln(os.Stderr, "[INFO] Reverse WHOIS succeeded but found no domains matching the extracted registrant details in current records.")
		}
	}
	if *wildcard {
		wildcardDomain := fmt.Sprintf("*%s*.%s", parts[0], parts[1])
		domains, err := core.GetDomains(wildcardDomain, *apiKey)
		if err != nil {
			fail(fmt.Errorf("enumeration failed; results may be incomplete: %w", err))
		}
		if len(domains) == 0 {
			fmt.Fprintf(os.Stderr, "[INFO] Domains & Subdomains Discovery succeeded but found no domains matching %q in the provider's database.\n", wildcardDomain)
		}
		outputDomains = append(outputDomains, domains...)
	}
	sort.Strings(outputDomains)
	unique := outputDomains[:0]
	for _, domain := range outputDomains {
		if len(unique) == 0 || unique[len(unique)-1] != domain {
			unique = append(unique, domain)
		}
	}
	outputDomains = unique
	if len(outputDomains) == 0 {
		fmt.Fprintln(os.Stderr, "[INFO] No domains found. See the lookup explanations above; an empty result does not mean the API key is invalid.")
	} else {
		fmt.Fprintf(os.Stderr, "[OK] Enumeration completed: %d unique domains found.\n", len(outputDomains))
	}
	if *output != "" {
		if err := core.OutputFile(outputDomains, *output); err != nil {
			fail(err)
		}
		fmt.Fprintf(os.Stderr, "[OK] Domains saved in %s\n", *output)
	} else {
		for _, domain := range outputDomains {
			fmt.Println(domain)
		}
	}
}
