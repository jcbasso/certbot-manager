package flags

import (
	"errors"

	"certbot-manager/internal/config"
)

// --- Email Flag ---

type EmailFlag struct{}

func init() { Register(&EmailFlag{}) }

// GenerateArgs now takes config structs
func (f *EmailFlag) GenerateArgs(certCfg config.Certificate, globalCfg config.Globals) ([]string, error) {
	email := ResolveString(certCfg.Email, globalCfg.Email) // Use helper
	if email == "" {
		return nil, errors.New("email is required but was not resolved from configuration")
	}
	return []string{"--email", email}, nil
}

// --- Agree TOS Flag ---

type AgreeTosFlag struct{}

func init() { Register(&AgreeTosFlag{}) }

func (f *AgreeTosFlag) GenerateArgs(_ config.Certificate, _ config.Globals) ([]string, error) {
	//	Automation requires true. Forcing '--agree-tos'
	return []string{"--agree-tos"}, nil
}

// --- Staging Flag ---

type StagingFlag struct{}

func init() { Register(&StagingFlag{}) }

func (f *StagingFlag) GenerateArgs(certCfg config.Certificate, globalCfg config.Globals) ([]string, error) {
	isStaging := ResolveBoolPtr(certCfg.Staging, globalCfg.Staging) // Use helper
	if isStaging {
		return []string{"--staging"}, nil
	}
	return nil, nil
}

// --- No EFF Email Flag ---

type NoEffEmailFlag struct{}

func init() { Register(&NoEffEmailFlag{}) }

func (f *NoEffEmailFlag) GenerateArgs(certCfg config.Certificate, globalCfg config.Globals) ([]string, error) {
	noEffEmail := ResolveBoolPtr(nil, globalCfg.NoEffEmail)
	if noEffEmail {
		return []string{"--no-eff-email"}, nil
	}
	return nil, nil
}

// --- Key Type Flag ---

type KeyTypeFlag struct{}

func init() { Register(&KeyTypeFlag{}) }

func (f *KeyTypeFlag) GenerateArgs(certCfg config.Certificate, globalCfg config.Globals) ([]string, error) {
	keyType := ResolveString(certCfg.KeyType, globalCfg.KeyType) // Use helper
	if keyType != "" {
		return []string{"--key-type", keyType}, nil
	}
	return nil, nil
}
