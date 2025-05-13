package monitor

import (
	"fmt"
	"time"

	"github.com/oliashish/vofo/internals/config"
	"github.com/oliashish/vofo/logger"
	"github.com/shirou/gopsutil/disk"
)

// Disk monitors disk usage for the root partition.
func Disk() error {
	log, err := logger.NewLogger()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	cfg := config.GetConfig()
	if cfg == nil {
		log.Error("Configuration not initialized")
		return fmt.Errorf("configuration not initialized")
	}

	log.Info("Starting Disk monitoring")
	interval := time.Duration(cfg.Interval * float64(time.Second))
	alertThreshold := time.Duration(cfg.AlertThreshold * float64(time.Second))

	// Track duration of threshold breach
	breachStart := time.Time{}
	breachActive := false

	for {
		usage, err := disk.Usage("/")
		if err != nil {
			log.Error(fmt.Sprintf("Failed to get Disk usage: %s", err))
			return fmt.Errorf("disk usage: %w", err)
		}

		usagePercent := usage.UsedPercent
		log.Info(fmt.Sprintf("Disk Usage: %.2f%% (Used: %d GB, Total: %d GB)",
			usagePercent, usage.Used/1024/1024/1024, usage.Total/1024/1024/1024))

		// Check threshold
		if usagePercent > cfg.DiskThreshold {
			if !breachActive {
				breachStart = time.Now()
				breachActive = true
			}
			breachDuration := time.Since(breachStart)
			log.Warn(fmt.Sprintf("Disk usage (%.2f%%) exceeds threshold (%.2f%%) for %.2f seconds",
				usagePercent, cfg.DiskThreshold, breachDuration.Seconds()))

			if breachDuration >= alertThreshold {
				message := fmt.Sprintf("Alert: Disk usage (%.2f%%) exceeded threshold (%.2f%%) for %.2f seconds",
					usagePercent, cfg.DiskThreshold, breachDuration.Seconds())
				if err := sendAlert(message, cfg, log); err != nil {
					log.Error(fmt.Sprintf("Failed to send alert: %s", err))
				} else {
					log.Info("Disk alert sent successfully")
				}
				breachActive = false
			}
		} else {
			if breachActive {
				log.Info("Disk usage returned to normal")
				breachActive = false
			}
		}

		time.Sleep(interval)
	}
}
