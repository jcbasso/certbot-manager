package authenticators

import (
	"fmt"
	"strconv"

	"certbot-manager/internal/certbot/flags" // Import flags for helpers
	"certbot-manager/internal/config"
)

type DuckDNSAuthenticator struct{}

func init() { Register("dns-duckdns", &DuckDNSAuthenticator{}) }

func (p *DuckDNSAuthenticator) BuildArgs(certCfg config.Certificate, globalCfg config.Globals) ([]string, error) {
	token := flags.ResolveString(certCfg.DuckDNSToken, globalCfg.DuckDNSToken)
	if token == "" {
		return nil, fmt.Errorf("authenticator 'dns-duckdns' requires the duckdns_token to be specified")
	}

	args := []string{
		"--authenticator", "dns-duckdns",
		"--dns-duckdns-token", token,
	}

	// Use helper to resolve propagation seconds
	propagationSecondsPtr := flags.ResolveIntPtr(certCfg.DNSPropagationSeconds, globalCfg.DNSPropagationSeconds)
	if propagationSecondsPtr == nil {
		return nil, fmt.Errorf("authenticator 'dns-duckdns' requires the dns_propagation_seconds to be specified")
	}

	propagationSeconds := *propagationSecondsPtr
	if propagationSeconds > 0 {
		args = append(args, "--dns-duckdns-propagation-seconds", strconv.Itoa(propagationSeconds))
	}

	return args, nil
}
