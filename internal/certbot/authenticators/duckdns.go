package authenticators

import (
	"fmt"
	"os"
	"strconv"

	"certbot-manager/internal/certbot/flags" // Import flags for helpers
	"certbot-manager/internal/config"
)

type DuckDNSAuthenticator struct{}

func init() { Register("dns-duckdns", &DuckDNSAuthenticator{}) }

func (p *DuckDNSAuthenticator) BuildArgs(certCfg config.Certificate, globalCfg config.Globals) ([]string, error) {
	token := os.Getenv("DUCKDNS_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("authenticator 'dns-duckdns' requires the DUCKDNS_TOKEN environment variable")
	}

	args := []string{
		"--authenticator", "dns-duckdns",
		"--dns-duckdns-token", token,
	}

	// Use helper to resolve propagation seconds
	propagationSeconds := flags.ResolveIntPtr(certCfg.DNSPropagationSeconds, config.Defaults.DNSPropagationSeconds)

	if propagationSeconds > 0 {
		args = append(args, "--dns-duckdns-propagation-seconds", strconv.Itoa(propagationSeconds))
	}

	return args, nil
}
