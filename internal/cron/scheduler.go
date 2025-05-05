package cron

import (
	"certbot-manager/internal/logging"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// Scheduler wraps the cron instance.
type Scheduler struct {
	instance *cron.Cron
}

// SetupAndStartScheduler initializes the cron scheduler, adds the specified job, and starts it.
// It takes the cron expression string and the job function to execute.
func SetupAndStartScheduler(expression string, job func()) (*Scheduler, error) {
	if expression == "" {
		return nil, fmt.Errorf("cron expression cannot be empty")
	}

	logrus.Info("--- Setting up Cron Scheduler ---")

	cronStdLogger := logging.NewLogrusStandardLogger(logrus.InfoLevel, "cron")
	cronLogger := cron.PrintfLogger(cronStdLogger)
	c := cron.New(
		cron.WithChain(
			cron.SkipIfStillRunning(cronLogger),
			cron.Recover(cronLogger),
		),
		cron.WithLogger(cronLogger),
		cron.WithSeconds(),
	)

	logrus.Infof("Scheduling job with cron expression: %s", expression)

	// Add the provided job function with the given expression
	entryID, err := c.AddFunc(expression, job)
	if err != nil {
		logrus.Errorf("Failed to add job to cron scheduler (expression: '%s'): %v", expression, err)
		return nil, fmt.Errorf("failed to add job to cron scheduler (expression: '%s'): %w", expression, err)
	}
	logrus.Infof("Renewal job added with ID: %d", entryID)

	c.Start()
	logrus.Info("Cron scheduler started.")

	return &Scheduler{instance: c}, nil
}

// Stop gracefully stops the cron scheduler, waiting for running jobs to complete.
func (s *Scheduler) Stop() {
	if s.instance != nil {
		logrus.Info("Stopping cron scheduler gracefully...")
		ctx := s.instance.Stop()
		<-ctx.Done()
		logrus.Info("Cron scheduler stopped.")
	} else {
		logrus.Warn("Scheduler instance is nil, cannot stop.")
	}
}
