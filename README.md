<h4 align="center">A fast tool for performing horizontal domain enumeration.</h4>

<p align="center">
  <a href="#installation-instructions">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#running-seekly">Running Seekly</a>
</p>

---

Seekly is a horizontal enumeration tool. This means the tool searches for domains related to an initial domain by leveraging [WhoisXMLAPI's](https://www.whoisxmlapi.com/) APIs.

# Installation Instructions
Seekly requires `go 1.24.1` to install successfully.
```sh
go install github.com/LucasKatashi/seekly/cmd/seekly@latest
```

# Usage
```sh
seekly -h
```

This will display help for the tool. Here are all the switches it supports.
```console
Usage:
 ./seekly [flags]

INPUT:
 --domain		enter the target domain
 --wildcard		performs a wildcard search to return all domains containing the value specified with --domain.
            for instance, using `--domain example.com` will match domains like *example*.com
 --api-key		enter your WhoisXMLAPI API key

OUTPUT:
 --output		generates an output file containing the discovered domains
 --silent		ignore the banner when running the tool
```

## Running Seekly
Seekly requires a WhoisXMLAPI API key to run. You can get your key [here](https://user.whoisxmlapi.com/products).

You can provide the key either by using the `--api-key`/`-api-key` flag or by setting it as an environment variable in your current session:
```sh
export WhoisXMLAPIKey="API_KEY"
```

Running Seekly **without** the `--wildcard` flag ensures you're only retrieving domains that are *actually related* to the one specified with `--domain`. By default, Seekly performs recursive lookups using WHOIS-related fields like email, person, owner, etc., to discover associated domains.

When using the `--wildcard` flag, Seekly instead performs a wildcard-based search for **any** domain containing the given domain value, restricted to the same TLD. For example:
```sh
--domain example.com --wildcard
```

Will result in a query like:
```sh
*example*.com
```

Seekly also integrates seamlessly in toolchains to further refine your recon. For example:
```sh
seekly -domain example.com -wildcard -silent | subfinder -all -silent | httpx -fr -mc 200 -t 150 -silent | katana -jc -jsl -d 4 -c 20 -silent --output example.txt
```

> ⚠️ Keep in mind that Seekly can significantly consume your WhoisXMLAPI credits, as it performs recursive lookups using WHOIS queries, Reverse WHOIS, the Domains & Subdomains API, and more.
