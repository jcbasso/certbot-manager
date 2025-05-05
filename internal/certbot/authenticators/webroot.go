package authenticators

import (
	"certbot-manager/internal/config"
)

type WebrootAuthenticator struct{}

func init() { Register("webroot", &WebrootAuthenticator{}) }

func (p *WebrootAuthenticator) BuildArgs(certCfg config.Certificate, globalCfg config.Globals) ([]string, error) {
	webrootPath := certCfg.WebrootPath
	if webrootPath == "" {
		webrootPath = globalCfg.WebrootPath // Fallback to global
	}

	return []string{"--webroot", "-w", webrootPath}, nil
}
