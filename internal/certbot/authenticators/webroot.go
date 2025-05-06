package authenticators

import (
	"certbot-manager/internal/certbot/flags"
	"certbot-manager/internal/config"
	"fmt"
)

type WebrootAuthenticator struct{}

func init() { Register("webroot", &WebrootAuthenticator{}) }

func (p *WebrootAuthenticator) BuildArgs(certCfg config.Certificate, globalCfg config.Globals) ([]string, error) {
	webrootPath := flags.ResolveString(certCfg.WebrootPath, globalCfg.WebrootPath)
	if webrootPath == "" {
		return nil, fmt.Errorf("authenticator 'webroot' requires the webroot_path to be specified")
	}

	return []string{"--webroot", "-w", webrootPath}, nil
}
