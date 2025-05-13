package monitor

import (
	"fmt"
	"time"

	"github.com/oliashish/vofo/internals/config"
	"github.com/oliashish/vofo/logger"
	"github.com/shirou/gopsutil/cpu"
)

// CPU monitors CPU usage for all cores.
func CPU() error {
	log, err := logger.NewLogger()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	cfg := config.GetConfig()
	if cfg == nil {
		log.Error("Configuration not initialized")
		return fmt.Errorf("configuration not initialized")
	}

	log.Info("Starting CPU monitoring")
	interval := time.Duration(cfg.Interval * float64(time.Second))
	alertThreshold := time.Duration(cfg.AlertThreshold * float64(time.Second))

	// Track duration of threshold breach
	breachStart := time.Time{}
	breachActive := false

	for {
		percentages, err := cpu.Percent(interval, true) // true = per-core
		if err != nil {
			log.Error(fmt.Sprintf("Failed to get CPU usage: %s", err))
			return fmt.Errorf("cpu usage: %w", err)
		}

		// Calculate average usage across all cores
		var total float64
		for i, perc := range percentages {
			total += perc
			log.Info(fmt.Sprintf("CPU Core %d: %.2f%%", i, perc))
		}
		avgUsage := total / float64(len(percentages))
		log.Info(fmt.Sprintf("Average CPU Usage: %.2f%%", avgUsage))

		// Check threshold
		if avgUsage > cfg.CPUThreshold {
			if !breachActive {
				breachStart = time.Now()
				breachActive = true
			}
			breachDuration := time.Since(breachStart)
			log.Warn(fmt.Sprintf("CPU usage (%.2f%%) exceeds threshold (%.2f%%) for %.2f seconds",
				avgUsage, cfg.CPUThreshold, breachDuration.Seconds()))

			if breachDuration >= alertThreshold {
				message := fmt.Sprintf("Alert: CPU usage (%.2f%%) exceeded threshold (%.2f%%) for %.2f seconds",
					avgUsage, cfg.CPUThreshold, breachDuration.Seconds())
				if err := sendAlert(message, cfg, log); err != nil {
					log.Error(fmt.Sprintf("Failed to send alert: %s", err))
				} else {
					log.Info("CPU alert sent successfully")
				}
				// Reset breach to avoid repeated alerts
				breachActive = false
			}
		} else {
			if breachActive {
				log.Info("CPU usage returned to normal")
				breachActive = false
			}
		}

		// Sleep for the interval (already accounted for in cpu.Percent)
	}
}

// sendAlert sends an alert based on the config's alert method (placeholder).
func sendAlert(message string, cfg *config.Config, log *logger.Logger) error {
	switch cfg.AlertMethod {
	case "email":
		// Placeholder: Implement email sending
		log.Info(fmt.Sprintf("Sending email to %s: %s", cfg.EmailRecipient, message))
		return nil
	case "slack":
		// Placeholder: Implement Slack webhook
		log.Info(fmt.Sprintf("Sending Slack message to %s: %s", cfg.SlackWebhook, message))
		return nil
	default:
		return fmt.Errorf("unsupported alert method: %s", cfg.AlertMethod)
	}
}
