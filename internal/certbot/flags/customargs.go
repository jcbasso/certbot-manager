package flags

import "certbot-manager/internal/config"

type CustomArgsFlag struct{}

func init() { Register(&CustomArgsFlag{}) }

func (f *CustomArgsFlag) GenerateArgs(certCfg config.Certificate, globalCfg config.Globals) ([]string, error) {
	if certCfg.Args != "" {
		return []string{certCfg.Args}, nil
	}
	return nil, nil
}
