package authenticators

import (
	"certbot-manager/internal/config"
)

// Authenticator defines the interface for different certbot challenge methods.
type Authenticator interface {
	// BuildArgs generates the certbot command line arguments specific to this authenticator.
	// It receives the specific certificate config for context.
	BuildArgs(certCfg config.Certificate, globalCfg config.Globals) ([]string, error)
}
