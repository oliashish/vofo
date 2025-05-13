package monitor

import (
	"fmt"
	"sort"
	"time"

	"github.com/oliashish/vofo/internals/config"
	"github.com/oliashish/vofo/logger"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

// Mem monitors RAM usage and top memory-consuming processes.
func Mem() error {
	log, err := logger.NewLogger()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	cfg := config.GetConfig()
	if cfg == nil {
		log.Error("Configuration not initialized")
		return fmt.Errorf("configuration not initialized")
	}

	log.Info("Starting RAM monitoring")
	interval := time.Duration(cfg.Interval * float64(time.Second))
	alertThreshold := time.Duration(cfg.AlertThreshold * float64(time.Second))

	// Track duration of threshold breach
	breachStart := time.Time{}
	breachActive := false

	for {
		vm, err := mem.VirtualMemory()
		if err != nil {
			log.Error(fmt.Sprintf("Failed to get RAM usage: %s", err))
			return fmt.Errorf("ram usage: %w", err)
		}

		usage := vm.UsedPercent
		log.Info(fmt.Sprintf("RAM Usage: %.2f%% (Used: %d MB, Total: %d MB)",
			usage, vm.Used/1024/1024, vm.Total/1024/1024))

		// Get top 3 memory-consuming processes
		processes, err := process.Processes()
		if err != nil {
			log.Warn(fmt.Sprintf("Failed to get processes: %s", err))
		} else {
			type procInfo struct {
				PID      int32
				Name     string
				MemUsage uint64
			}
			var procList []procInfo

			for _, p := range processes {
				mem, err := p.MemoryInfo()
				if err != nil {
					continue
				}
				name, err := p.Name()
				if err != nil {
					name = "unknown"
				}
				procList = append(procList, procInfo{
					PID:      p.Pid,
					Name:     name,
					MemUsage: mem.RSS / 1024 / 1024, // Convert to MB
				})
			}

			// Sort by memory usage (descending)
			sort.Slice(procList, func(i, j int) bool {
				return procList[i].MemUsage > procList[j].MemUsage
			})

			// Log top 3 processes
			maxProcs := 3
			if len(procList) < maxProcs {
				maxProcs = len(procList)
			}
			for i := 0; i < maxProcs; i++ {
				log.Info(fmt.Sprintf("Top Process %d: PID=%d, Name=%s, Mem=%d MB",
					i+1, procList[i].PID, procList[i].Name, procList[i].MemUsage))
			}
		}

		// Check threshold
		if usage > cfg.RAMThreshold {
			if !breachActive {
				breachStart = time.Now()
				breachActive = true
			}
			breachDuration := time.Since(breachStart)
			log.Warn(fmt.Sprintf("RAM usage (%.2f%%) exceeds threshold (%.2f%%) for %.2f seconds",
				usage, cfg.RAMThreshold, breachDuration.Seconds()))

			if breachDuration >= alertThreshold {
				message := fmt.Sprintf("Alert: RAM usage (%.2f%%) exceeded threshold (%.2f%%) for %.2f seconds",
					usage, cfg.RAMThreshold, breachDuration.Seconds())
				if err := sendAlert(message, cfg, log); err != nil {
					log.Error(fmt.Sprintf("Failed to send alert: %s", err))
				} else {
					log.Info("RAM alert sent successfully")
				}
				breachActive = false
			}
		} else {
			if breachActive {
				log.Info("RAM usage returned to normal")
				breachActive = false
			}
		}

		time.Sleep(interval)
	}
}
