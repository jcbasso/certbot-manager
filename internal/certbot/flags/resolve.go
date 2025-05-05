package flags

import (
	"certbot-manager/internal/config"
)

// --- Resolution Helper Functions ---

func ResolveString(certVal, globalVal string) string {
	if certVal != "" {
		return certVal
	}
	return globalVal
}

func ResolveBoolPtr(certVal *bool, globalVal bool) bool {
	if certVal != nil {
		return *certVal
	}
	return globalVal
}

// ResolveIntPtr handles integer overrides with a default fallback
func ResolveIntPtr(certVal *int, defaultVal int) int {
	if certVal != nil {
		return *certVal
	}
	return defaultVal
}

func ResolveAuthenticatorName(certCfg config.Certificate, globalCfg config.Globals) string {
	// Authenticator is primarily per-certificate
	if certCfg.Authenticator != "" {
		return certCfg.Authenticator
	}
	// Fallback to a hardcoded default if not specified
	return config.DefaultAuthenticator
}
