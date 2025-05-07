package authenticators

import (
	"fmt"
	"strconv"

	"certbot-manager/internal/certbot/flags"
	"certbot-manager/internal/config"
)

type CloudflareAuthenticator struct{}

func init() { Register("dns-cloudflare", &CloudflareAuthenticator{}) }

func (p *CloudflareAuthenticator) BuildArgs(certCfg config.Certificate, globalCfg config.Globals) ([]string, error) {
	credentialsPath := flags.ResolveString(certCfg.CloudflareCredentialsPath, globalCfg.CloudflareCredentialsPath)

	if credentialsPath == "" {
		return nil, fmt.Errorf("authenticator 'dns-cloudflare' requires the cloudflare_credentials_path to be specified")
	}

	// Base arguments for the plugin
	args := []string{
		"--authenticator", "dns-cloudflare",
		"--dns-cloudflare-credentials", credentialsPath,
	}

	// Use helper to resolve propagation seconds
	propagationSecondsPtr := flags.ResolveIntPtr(certCfg.DNSPropagationSeconds, globalCfg.DNSPropagationSeconds)
	if propagationSecondsPtr == nil {
		return nil, fmt.Errorf("authenticator 'dns-cloudflare' requires the dns_propagation_seconds to be specified")
	}

	propagationSeconds := *propagationSecondsPtr
	if propagationSeconds > 0 {
		args = append(args, "--dns-cloudflare-propagation-seconds", strconv.Itoa(propagationSeconds))
	}

	return args, nil
}
