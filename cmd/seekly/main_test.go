package main

import (
	"flag"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

type responseTransport struct {
	whois   string
	domains string
}

func (f responseTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	body := f.domains
	switch r.URL.Path {
	case "/user-service/account-balance":
		body = `{"data":[{"product":{"name":"WHOIS API"},"credits":10},{"product":{"name":"Domain Research Suite"},"credits":10}]}`
	case "/whoisserver/WhoisService":
		body = f.whois
	}
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
}

func TestCLIExplainsEmptyResultsAndKeepsStdoutClean(t *testing.T) {
	for _, tc := range []struct {
		name, whois, domains, stdout, message string
	}{
		{"no registrant", `{"WhoisRecord":{"domainName":"example.com"}}`, `{"domainsList":[]}`, "", "Reverse WHOIS was skipped"},
		{"no matches", `{"WhoisRecord":{"registrant":{"name":"Example"}}}`, `{"domainsList":[]}`, "", "found no domains matching the extracted registrant details"},
		{"domains", `{"WhoisRecord":{"registrant":{"name":"Example"}}}`, `{"domainsList":["b.example","a.example","a.example"]}`, "a.example\nb.example\n", "2 unique domains found"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			oldArgs, oldFlags := os.Args, flag.CommandLine
			oldOut, oldErr, oldTransport := os.Stdout, os.Stderr, http.DefaultTransport
			out, err := os.CreateTemp(t.TempDir(), "stdout")
			if err != nil {
				t.Fatal(err)
			}
			diagnostics, err := os.CreateTemp(t.TempDir(), "stderr")
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				os.Args, flag.CommandLine = oldArgs, oldFlags
				os.Stdout, os.Stderr, http.DefaultTransport = oldOut, oldErr, oldTransport
				out.Close()
				diagnostics.Close()
			})
			t.Setenv("WhoisXMLAPIKey", "test-secret")
			os.Args = []string{"seekly", "--domain", "example.com", "--silent", "--wildcard"}
			flag.CommandLine = flag.NewFlagSet("seekly", flag.ContinueOnError)
			os.Stdout, os.Stderr = out, diagnostics
			http.DefaultTransport = responseTransport{tc.whois, tc.domains}
			main()
			stdout, err := os.ReadFile(out.Name())
			if err != nil {
				t.Fatal(err)
			}
			stderr, err := os.ReadFile(diagnostics.Name())
			if err != nil {
				t.Fatal(err)
			}
			if string(stdout) != tc.stdout {
				t.Fatalf("stdout=%q, want %q", stdout, tc.stdout)
			}
			if !strings.Contains(string(stderr), tc.message) {
				t.Fatalf("missing explanation %q: %s", tc.message, stderr)
			}
			if tc.stdout == "" && (!strings.Contains(string(stderr), "No domains found") || !strings.Contains(string(stderr), "Discovery succeeded but found no domains")) {
				t.Fatalf("missing empty-result explanation: %s", stderr)
			}
		})
	}
}
