package core

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type fixtureTransport struct {
	url string
}

func (f fixtureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	copy := req.Clone(req.Context())
	copy.URL.Scheme = "http"
	copy.URL.Host = strings.TrimPrefix(f.url, "http://")
	return http.DefaultTransport.RoundTrip(copy)
}

func fixture(t *testing.T, status int, body string) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}))
	previous := apiClient
	apiClient = &http.Client{Transport: fixtureTransport{server.URL}}
	t.Cleanup(func() { apiClient = previous; server.Close() })
}

func TestLookupFailuresAreNotEmptyResults(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		want   string
	}{
		{"access denied", 403, `{"code":403,"messages":"insufficient credits"}`, "insufficient credits"},
		{"unauthorized", 401, `{"messages":"incorrect API key"}`, "access denied"},
		{"rate limited", 429, `{"messages":"Too many requests"}`, "rate limit"},
		{"server failure", 502, `<html>Bad Gateway</html>`, "HTTP 502"},
		{"WHOIS error envelope", 200, `{"ErrorMessage":{"msg":"API key test-secret rejected"}}`, "[REDACTED]"},
		{"API error envelope", 200, `{"code":403,"messages":"insufficient credits"}`, "insufficient credits"},
		{"invalid JSON", 200, `<html>blocked</html>`, "invalid JSON"},
		{"missing result", 200, `{}`, "missing a valid"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fixture(t, tc.status, tc.body)
			for _, lookup := range []func() ([]string, error){
				func() ([]string, error) { return GetWhois("example.com", "test-secret") },
				func() ([]string, error) { return GetRevWhois([]string{"Example"}, "test-secret") },
				func() ([]string, error) { return GetDomains("*example*.com", "test-secret") },
			} {
				_, err := lookup()
				if err == nil || !strings.Contains(err.Error(), tc.want) || strings.Contains(err.Error(), "test-secret") {
					t.Fatalf("expected %q without secret, got %v", tc.want, err)
				}
			}
		})
	}
}

func TestValidEmptySearches(t *testing.T) {
	fixture(t, 200, `{"domainsList":[],"domainsCount":0}`)
	for _, lookup := range []func() ([]string, error){
		func() ([]string, error) { return GetRevWhois([]string{"Example"}, "test-secret") },
		func() ([]string, error) { return GetDomains("*example*.com", "test-secret") },
	} {
		domains, err := lookup()
		if err != nil || len(domains) != 0 {
			t.Fatalf("valid empty result: domains=%v error=%v", domains, err)
		}
	}
}

func TestReverseWhoisReturnsOnlyUniqueDomains(t *testing.T) {
	fixture(t, 200, `{"domainsList":["b.example","a.example","a.example"],"registrant":{"name":"not-a-domain"}}`)
	domains, err := GetRevWhois([]string{"Example", "Example Inc"}, "test-secret")
	if err != nil || !reflect.DeepEqual(domains, []string{"a.example", "b.example"}) {
		t.Fatalf("domains=%v error=%v", domains, err)
	}
}

func TestWhoisWithoutRegistrant(t *testing.T) {
	fixture(t, 200, `{"WhoisRecord":{"domainName":"example.com"}}`)
	values, err := GetWhois("example.com", "test-secret")
	if err != nil || len(values) != 0 {
		t.Fatalf("expected no registrant pivots, got %v, %v", values, err)
	}
}

func TestWhoisDataError(t *testing.T) {
	fixture(t, 200, `{"WhoisRecord":{"dataError":"MISSING_WHOIS_DATA"}}`)
	_, err := GetWhois("example.com", "test-secret")
	if err == nil || !strings.Contains(err.Error(), "MISSING_WHOIS_DATA") {
		t.Fatalf("expected explicit missing WHOIS data error, got %v", err)
	}
}

func TestProductBalances(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
		want string
	}{
		{"WHOIS cannot pay for DRS", `{"data":[{"product":{"name":"WHOIS API"},"credits":379},{"product":{"name":"Domain Research Suite"},"credits":0}]}`, "Domain Research Suite has no available credits"},
		{"WHOIS exhausted", `{"data":[{"product":{"name":"WHOIS API"},"credits":0},{"product":{"name":"Domain Research Suite"},"credits":10}]}`, "WHOIS API has no available credits"},
		{"missing product", `{"data":[]}`, "not present"},
		{"malformed balance", `{}`, "inconclusive"},
		{"unknown credits", `{"data":[{"product":{"name":"WHOIS API"}}]}`, "cannot determine"},
		{"funded", `{"data":[{"product":{"name":"WHOIS API"},"credits":10},{"product":{"name":"Domain Research Suite"},"credits":10}]}`, ""},
		{"overage enabled", `{"data":[{"product":{"name":"WHOIS API"},"credits":0,"overage_eligible":true},{"product":{"name":"Domain Research Suite"},"credits":0,"overage_eligible":true}]}`, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fixture(t, 200, tc.body)
			err := CheckAPIKey("test-secret")
			if tc.want == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected %q, got %v", tc.want, err)
			}
		})
	}
}
