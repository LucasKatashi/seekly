package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const balanceAPI = "https://user.whoisxmlapi.com/user-service/account-balance"

var apiClient = &http.Client{Timeout: 30 * time.Second}

// requestAPI distinguishes provider failures from valid empty search results.
func requestAPI(endpoint, product, apiKey string, payload any) (map[string]json.RawMessage, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("%s: cannot encode request: %w", product, err)
	}
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("%s: cannot create request: %w", product, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "seekly")
	resp, err := apiClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed; API key validity could not be determined: %s", product, strings.ReplaceAll(err.Error(), apiKey, "[REDACTED]"))
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: cannot read response: %w", product, err)
	}
	var object map[string]json.RawMessage
	decodeErr := json.Unmarshal(body, &object)
	var details []string
	for _, field := range []string{"ErrorMessage", "error", "errors", "messages"} {
		if value := object[field]; len(value) > 0 && string(value) != "null" && string(value) != `""` && string(value) != "[]" && string(value) != "{}" {
			details = append(details, string(value))
		}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 || len(details) > 0 {
		reason := strings.Join(details, "; ")
		if reason == "" {
			reason = "provider returned an unsuccessful response"
		}
		switch resp.StatusCode {
		case http.StatusUnauthorized, http.StatusForbidden:
			reason += "; access denied: check the API key, product credits/subscription, and IP allowlist. This does not by itself prove the key is invalid"
		case http.StatusTooManyRequests:
			reason += "; rate limit exceeded; try again later"
		}
		return nil, fmt.Errorf("%s (HTTP %d): %s", product, resp.StatusCode, strings.ReplaceAll(reason, apiKey, "[REDACTED]"))
	}
	if decodeErr != nil || object == nil {
		return nil, fmt.Errorf("%s (HTTP %d): invalid JSON response; this is not an empty search result", product, resp.StatusCode)
	}
	return object, nil
}

// CheckAPIKey checks authentication separately from the credits needed by Seekly.
func CheckAPIKey(apiKey string) error {
	object, err := requestAPI(balanceAPI, "Account balance API", apiKey, map[string]string{"apiKey": apiKey, "outputFormat": "JSON"})
	if err != nil {
		return fmt.Errorf("API key check failed: %w", err)
	}
	var balances []struct {
		Product struct {
			Name string `json:"name"`
		} `json:"product"`
		Credits         *float64 `json:"credits"`
		OverageEligible bool     `json:"overage_eligible"`
	}
	if err := json.Unmarshal(object["data"], &balances); err != nil || balances == nil {
		return fmt.Errorf("API key check inconclusive: account balance response is missing a valid data list")
	}
	fmt.Fprintln(os.Stderr, "[OK] API key is working: account balance authentication succeeded. Product access is checked separately.")
	var unavailable []string
	for _, name := range []string{"WHOIS API", "Domain Research Suite"} {
		found := false
		for _, balance := range balances {
			if balance.Product.Name != name {
				continue
			}
			found = true
			if balance.Credits == nil {
				return fmt.Errorf("cannot determine %s credits from the account balance response", name)
			}
			fmt.Fprintf(os.Stderr, "[INFO] %s credits: %g; overage eligible: %t.\n", name, *balance.Credits, balance.OverageEligible)
			if *balance.Credits <= 0 && !balance.OverageEligible {
				unavailable = append(unavailable, name+" has no available credits and overage is disabled")
			}
			break
		}
		if !found {
			unavailable = append(unavailable, name+" is not present in the account balance")
		}
	}
	if len(unavailable) > 0 {
		return fmt.Errorf("no domains returned: %s. Reverse WHOIS and Domains & Subdomains Discovery require Domain Research Suite credits; WHOIS API credits cannot replace them", strings.Join(unavailable, "; "))
	}
	return nil
}

func domainList(object map[string]json.RawMessage, product string) ([]string, error) {
	var domains []string
	if err := json.Unmarshal(object["domainsList"], &domains); err != nil || domains == nil {
		return nil, fmt.Errorf("%s: response is missing a valid domainsList; cannot treat it as zero matches", product)
	}
	return domains, nil
}
