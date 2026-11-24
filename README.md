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

Without `--wildcard`, Seekly queries WHOIS for the initial domain, then searches current Reverse WHOIS records using the registrant name, organization, and email fields it can extract. Shared or privacy-proxy contact details can produce unrelated matches; results are leads, not proof of common ownership. Lookups are not recursive.

With `--wildcard`, Seekly also performs a wildcard-based search for domains containing the given value, restricted to the suffix after the first dot, and merges these matches with the Reverse WHOIS results. For example:
```sh
--domain example.com --wildcard
```

Will result in a query like:
```sh
*example*.com
```

Seekly also integrates seamlessly in toolchains to further refine your recon. For example:
```sh
seekly -domain example.com -silent | subfinder -all -silent | httpx -fr -mc 200 -t 150 -silent | katana -jc -jsl -d 4 -c 20 -silent --output example.txt
```

Seekly consumes WHOIS API credits for the initial lookup and separate **Domain Research Suite (DRS)** credits for Reverse WHOIS and Domains & Subdomains Discovery. A working key with WHOIS credits but no DRS credits cannot complete enumeration.

### Authentication and empty results

Before enumeration, Seekly checks the [account balance endpoint](https://whois.whoisxmlapi.com/documentation/account-balance-information). It reports successful authentication separately from product balances, and stops before paid lookups if a required product is missing or has no credits and overage is disabled.

- Authentication success does not guarantee access to every product. HTTP 401/403 can indicate a restricted key, insufficient credits, or an IP allowlist restriction; Seekly includes the provider's error without claiming that every denied key is invalid.
- HTTP failures, rate limits, timeouts, API error envelopes, invalid JSON, and unexpected response formats produce explicit errors, not empty search results.
- Missing WHOIS data is reported separately. A valid WHOIS record with no usable registrant fields explains why Reverse WHOIS was skipped; privacy/redaction or absent fields may be responsible.
- Successful searches with zero matches explain which lookup found nothing. No matches in the provider's database do not prove that no related domains exist.

Diagnostics and the banner go to **stderr**; stdout contains only discovered domains. `--silent` hides only the banner, not authentication, balance, warning, or error messages. Errors exit with status **1**; successful searches, including genuinely empty results, exit with status **0**. Enumeration failures do not write a misleading empty output file. A successful empty search with `--output` writes an empty file.

To install local changes:

```sh
go install ./cmd/seekly
```
