package certbot

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"certbot-manager/internal/config" // Import config package
)

// ValidateCertbotPath checks if the certbot command exists and is executable.
// Returns the absolute path to the executable or an error.
func ValidateCertbotPath(potentialPath string) (string, error) {
	if potentialPath == "" {
		return "", errors.New("certbot path configuration is empty")
	}

	resolvedPath, err := exec.LookPath(potentialPath)
	if err != nil {
		if _, statErr := os.Stat(potentialPath); os.IsNotExist(statErr) {
			return "", fmt.Errorf("certbot executable '%s' not found in PATH and does not exist: %w", potentialPath, err)
		} else if statErr != nil {
			return "", fmt.Errorf("error checking certbot path '%s': %w", potentialPath, statErr)
		}
		return "", fmt.Errorf("certbot executable '%s' not found in PATH: %w", potentialPath, err)
	}

	fileInfo, err := os.Stat(resolvedPath)
	if err != nil {
		return "", fmt.Errorf("could not stat resolved certbot path '%s': %w", resolvedPath, err)
	}
	if fileInfo.Mode()&0100 == 0 { // Check user execute permission bit
		return "", fmt.Errorf("resolved certbot path '%s' is not executable", resolvedPath)
	}

	logrus.Infof("Validated certbot executable: %s", resolvedPath)
	return resolvedPath, nil
}

// runCommand executes the certbot command with given arguments.
func runCommand(executablePath string, args ...string) error {
	cmd := exec.Command(executablePath, args...)
	logrus.Debugf("Running command: %s %s", executablePath, strings.Join(args, " "))

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run() // Waits for completion

	stdoutStr := strings.TrimSpace(stdoutBuf.String())
	stderrStr := strings.TrimSpace(stderrBuf.String())

	if len(stdoutStr) > 0 {
		logrus.Debugf("Command stdout:\n---\n%s\n---", stdoutStr)
	}

	if err != nil {
		var exitErr *exec.ExitError
		exitCode := -1
		if errors.As(err, &exitErr) {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}
		errMsg := fmt.Sprintf("Command failed with error: %v", err)
		if len(stderrStr) > 0 {
			errMsg += fmt.Sprintf("\nStderr:\n---\n%s\n---", stderrStr)
		}
		logrus.Errorf("%s (Exit Code: %d)", errMsg, exitCode)
		return fmt.Errorf("command execution failed (exit code %d): %w", exitCode, err)
	}

	logrus.Infof("Command finished successfully (Exit Code: 0)")
	return nil
}

// RequestCertificates handles the initial 'certbot certonly' runs for all configured certificates.
func RequestCertificates(cfg *config.Config, certbotPath string) bool { // Accepts *config.Config
	logrus.Info("--- Initial Certificate Processing ---")
	allRunsSuccessful := true

	for i, cert := range cfg.Certificates {
		logrus.Infof("Processing certificate request %d for domains: %v", i+1, cert.Domains)

		// Create builder with specific cert config and global config
		builder := NewArgsBuilder(cert, cfg.Globals)
		args, err := builder.Build()
		if err != nil {
			logrus.Errorf("Error building arguments for cert #%d (%v): %v. Skipping.", i+1, cert.Domains, err)
			allRunsSuccessful = false
			continue
		}

		err = runCommand(certbotPath, args...)
		if err != nil {
			logrus.Errorf("Failed initial certonly run for cert %d (%v): %v", i+1, cert.Domains, err)
			allRunsSuccessful = false
		}
	}
	return allRunsSuccessful
}

// RenewCertificates runs 'certbot renew'.
func RenewCertificates(certbotPath string) error {
	logrus.Info("Checking for certificate renewals...")
	err := runCommand(certbotPath, "renew", "--quiet")
	if err != nil {
		logrus.Infof("Certbot renew command finished with potential issue: %v", err)
		return err
	}
	logrus.Info("Certbot renew command finished.")
	return nil
}
