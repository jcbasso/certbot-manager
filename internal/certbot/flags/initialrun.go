package flags

import (
	"certbot-manager/internal/config" // Use correct module path
)

// InitialRunFlags handles --force-renewal or --keep-until-expiring.
type InitialRunFlags struct{}

func init() { Register(&InitialRunFlags{}) } // Register this generator

// GenerateArgs now takes config structs
func (f *InitialRunFlags) GenerateArgs(certCfg config.Certificate, globalCfg config.Globals) ([]string, error) {
	force := ResolveBoolPtr(certCfg.InitialForceRenewal, globalCfg.InitialForceRenewal) // Use helper
	if force != nil && *force {
		return []string{"--force-renewal"}, nil
	}
	return []string{"--keep-until-expiring"}, nil
}
