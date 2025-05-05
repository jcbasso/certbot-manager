package certbot

import (
	"certbot-manager/internal/certbot/authenticators"
	"certbot-manager/internal/certbot/flags" // Import flags package
	"certbot-manager/internal/config"
	"errors"
	"fmt"
)

// ArgsBuilder now primarily holds the context (configs) needed during the build process.
// It no longer holds resolved state itself.
type ArgsBuilder struct {
	certCfg   config.Certificate
	globalCfg config.Globals
}

// NewArgsBuilder simply stores the configuration context.
func NewArgsBuilder(cfg config.Certificate, globals config.Globals) *ArgsBuilder {
	return &ArgsBuilder{
		certCfg:   cfg,
		globalCfg: globals,
	}
}

// Build constructs the final argument list using flag generators and the authenticator plugin,
// passing the config context to them.
func (b *ArgsBuilder) Build() ([]string, error) {
	// Minimal validation here, more specific validation happens in generators/plugins
	if len(b.certCfg.Domains) == 0 {
		return nil, errors.New("at least one domain is required")
	}

	// Base Command
	args := []string{"certonly", "--non-interactive"}

	// Apply Common Flags via Registered Generators
	for _, generator := range flags.GetAll() {
		// Pass the config context directly to the generator
		flagArgs, err := generator.GenerateArgs(b.certCfg, b.globalCfg)
		if err != nil {
			typeName := fmt.Sprintf("%T", generator)
			return nil, fmt.Errorf("error from flag generator %s for domains %v: %w", typeName, b.certCfg.Domains, err)
		}
		if len(flagArgs) > 0 {
			args = append(args, flagArgs...)
		}
	}

	// Get and Apply Authenticator Args
	authenticatorName := flags.ResolveAuthenticatorName(b.certCfg, b.globalCfg) // Use helper for consistency
	plugin, err := authenticators.Get(authenticatorName)
	if err != nil {
		return nil, fmt.Errorf("failed to get authenticator plugin for '%s' (domains: %v): %w", authenticatorName, b.certCfg.Domains, err)
	}
	// Authenticator plugin interface already expects configs
	authArgs, err := plugin.BuildArgs(b.certCfg, b.globalCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build args for authenticator '%s' (domains: %v): %w", authenticatorName, b.certCfg.Domains, err)
	}
	if len(authArgs) > 0 {
		args = append(args, authArgs...)
	}

	for _, domain := range b.certCfg.Domains {
		args = append(args, "-d", domain)
	}

	return args, nil
}
