package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/LucasKatashi/seekly/internal/core"
	"github.com/LucasKatashi/seekly/internal/ui"
	"github.com/fatih/color"
)

func main() {
	domain := flag.String("domain", "", "")
	silent := flag.Bool("silent", false, "")
	wildcard := flag.Bool("wildcard", false, "")
	apiKey := flag.String("api-key", "", "")
	output := flag.String("output", "", "")

	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)

	flag.Usage = ui.CustomUsage
	flag.Parse()
	ui.PrintBanner(*silent)

	if *apiKey == "" {
		*apiKey = os.Getenv("WhoisXMLAPIKey")
		if *apiKey == "" {
			fmt.Printf("[%s] The \"WhoisXMLAPIKey\" API key is not defined either in your environment variables or via a flag.\n      To do so, use `export WhoisXMLAPIKey=\"API_KEY\"` or `--api-key API_KEY` when running the tool.\n\n      The API key can be obtained at: https://user.whoisxmlapi.com/products\n", red.Sprint("ERR"))
			os.Exit(1)
		}
	}

	if *domain != "" {
		values := core.GetWhois(*domain, *apiKey)

		whoisDomains := core.GetRevWhois(values, *apiKey)

		var outputDomains []string

		if *wildcard {
			parts := strings.SplitN(*domain, ".", 2)
			wildcardDomain := fmt.Sprintf("*%s*.%s", parts[0], parts[1])

			domains := core.GetDomains(wildcardDomain, *apiKey)

			allDomains := append(whoisDomains, domains...)

			uniqueDomains := make(map[string]bool)
			for _, d := range allDomains {
				uniqueDomains[d] = true
			}

			for domain := range uniqueDomains {
				outputDomains = append(outputDomains, domain)
			}
		} else {
			outputDomains = whoisDomains
		}

		if *output != "" {
			core.OutputFile(outputDomains, *output)
			fmt.Printf("[%s] Domains saved in %s\n", green.Sprint("OK"), *output)
		} else {
			for _, domain := range outputDomains {
				fmt.Println(domain)
			}
		}
	}
}
