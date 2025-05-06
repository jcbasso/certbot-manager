package flags

import (
	"certbot-manager/internal/config"
	"fmt"
)

// --- Resolution Helper Functions ---

func ResolveString(certVal, globalVal string) string {
	if certVal != "" {
		return certVal
	}
	return globalVal
}

func ResolveBoolPtr(certVal *bool, globalVal *bool) *bool {
	if certVal != nil {
		return certVal
	}
	return globalVal
}

// ResolveIntPtr handles integer overrides with a default fallback
func ResolveIntPtr(certVal *int, globalVal *int) *int {
	if certVal != nil {
		return certVal
	}
	return globalVal
}

func ResolveAuthenticatorName(certCfg config.Certificate, globalCfg config.Globals) (string, error) {
	authenticator := ResolveString(certCfg.Authenticator, globalCfg.Authenticator)
	if authenticator != "" {
		return authenticator, nil
	}

	return "", fmt.Errorf("no authenticator specified in config")
}
