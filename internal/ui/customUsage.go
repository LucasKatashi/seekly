package ui

import (
	"fmt"
	"os"
)

func CustomUsage() {
	fmt.Fprintf(os.Stderr, "Seekly is a horizontal enumeration tool. This means the tool searches for domains related to an initial domain by leveraging WhoisXMLAPI's APIs.\n\n")

	fmt.Fprintf(os.Stderr, "Usage:\n seekly [flags]\n\n")
	fmt.Fprintf(os.Stderr, "INPUT:\n")
	fmt.Fprintf(os.Stderr, " --domain\t\tenter the target domain\n")
	fmt.Fprintf(os.Stderr, " --wildcard\t\tperforms a wildcard search to return all domains containing the value specified with --domain.\n\t\t\tfor instance, using `--domain example.com` will match domains like *example*.com\n")
	fmt.Fprintf(os.Stderr, " --api-key\t\tenter your WhoisXMLAPI API key\n\n")

	fmt.Fprintf(os.Stderr, "OUTPUT:\n")
	fmt.Fprintf(os.Stderr, " --output\t\tgenerates an output file containing the discovered domains\n")
	fmt.Fprintf(os.Stderr, " --silent\t\tignore the banner when running the tool\n")
}
