package certbot

import (
	"certbot-manager/internal/certbot/flags"
	"certbot-manager/internal/config"
	"fmt"
)

var commandList = []string{
	"certonly",
	"run",
}

func generateCmd(certCfg config.Certificate, globalCfg config.Globals) (string, error) {
	cmd := flags.ResolveString(certCfg.Cmd, globalCfg.Cmd) // Use helper
	if cmd == "" {
		return config.Defaults.Cmd, nil
	}

	if !isValidCommand(cmd) {
		return "", fmt.Errorf("unknown cmd '%s' (options: 'certonly', 'run')", cmd)
	}
	return cmd, nil
}

func isValidCommand(cmd string) bool {
	for _, validCmd := range commandList {
		if cmd == validCmd {
			return true
		}
	}
	return false
}
