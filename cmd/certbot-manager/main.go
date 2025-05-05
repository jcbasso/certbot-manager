package main

import (
	"certbot-manager/internal/logging"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"

	"certbot-manager/internal/certbot"
	"certbot-manager/internal/config"
	cronpkg "certbot-manager/internal/cron"
)

func main() {
	// --- Load Configuration ---
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// --- Setup Logging ---
	if err := logging.Setup(cfg.LogLevel); err != nil {
		log.Fatalf("Failed to setup logging: %v", err)
	}

	logrus.Info("Starting Certbot Manager...")

	// --- Validate Certbot Path ---
	validatedCertbotPath, err := certbot.ValidateCertbotPath(cfg.CertbotPath)
	if err != nil {
		logrus.Fatalf("Certbot path validation failed: %v", err)
	}

	// --- Check if certificates need processing ---
	if len(cfg.Certificates) == 0 {
		logrus.Info("No [[certificate]] blocks found in configuration. Nothing to schedule.")
		os.Exit(0)
	}

	// --- Initial Certificate Request ---
	initialRunsOk := certbot.RequestCertificates(cfg, validatedCertbotPath)

	// --- !!! Check for Initial Failures !!! ---
	if !initialRunsOk {
		logrus.Fatal("FATAL: One or more initial certificate requests failed. " +
			"Check logs above for details. Application will not start the renewal scheduler.",
		)
	}

	// --- Proceed only if initial requests were successful ---
	logrus.Info("Initial certificates processing completed successfully.")

	// --- Define the Renewal Job Function ---
	renewalJob := func() {
		logrus.Info("Cron Job: Triggered renewal check...")
		err := certbot.RenewCertificates(validatedCertbotPath)
		if err != nil {
			logrus.Warn("Cron Job: Renewal check finished with potential issue.")
		} else {
			logrus.Info("Cron Job: Renewal check finished successfully.")
		}
	}

	// --- Setup and Start Cron Scheduler ---
	scheduler, err := cronpkg.SetupAndStartScheduler(cfg.Globals.RenewalCron, renewalJob)
	if err != nil {
		logrus.Fatalf("Failed to setup and start cron scheduler: %v", err)
	}

	// --- Wait for Shutdown Signal ---
	logrus.Info("Certbot Manager running. Renewal checks scheduled via cron. Waiting for signals...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// --- Initiate Graceful Shutdown ---
	logrus.Info("Shutdown signal received...")
	scheduler.Stop()

	logrus.Info("Certbot Manager application stopped.")
}
