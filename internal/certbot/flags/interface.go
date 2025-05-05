package flags

import (
	"certbot-manager/internal/config"
)

// FlagGenerator defines the interface for generating specific Certbot CLI flags.
type FlagGenerator interface {
	// GenerateArgs produces the command line arguments for this flag based directly
	// on the certificate and global configuration.
	// Returns the arguments slice (can be nil/empty if flag not applicable) and an optional error.
	GenerateArgs(certCfg config.Certificate, globalCfg config.Globals) ([]string, error)
}
