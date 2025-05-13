package cmd

import (
	"fmt"
	"strings"

	"github.com/oliashish/vofo/internals/config"
	"github.com/oliashish/vofo/internals/monitor"
)

// startAgent starts the agent based on the specified state.
func startAgent(state State) error {
	switch state {
	case All:
		return monitor.All()
	case Monitor:
		return monitor.All()
	case Heal:
		return monitor.Heal()
	default:
		log.Info("Defaulting to monitoring module")
		return monitor.All()
	}
}

// startMonitor starts monitoring the specified resources.
func startMonitor(resources []Resource) error {
	cfg := config.GetConfig()
	if cfg == nil {
		log.Error("No config loaded; run 'vofo init' first")
		return fmt.Errorf("configuration not initialized")
	}
	log.Info(fmt.Sprintf("Using thresholds: CPU=%.2f, RAM=%.2f, Disk=%.2f",
		cfg.CPUThreshold, cfg.RAMThreshold, cfg.DiskThreshold))

	if len(resources) == 0 {
		return monitor.All()
	}

	for _, r := range resources {
		switch r {
		case CPU:
			if err := monitor.CPU(); err != nil {
				log.Error(fmt.Sprintf("Failed to monitor CPU: %s", err))
				return err
			}
		case Disk:
			if err := monitor.Disk(); err != nil {
				log.Error(fmt.Sprintf("Failed to monitor Disk: %s", err))
				return err
			}
		case Mem:
			if err := monitor.Mem(); err != nil {
				log.Error(fmt.Sprintf("Failed to monitor Memory: %s", err))
				return err
			}
		}
	}
	return nil
}

// parseResources validates and parses resource arguments.
func parseResources(args []string, allFlag bool) ([]Resource, error) {
	if allFlag {
		return []Resource{CPU, Disk, Mem}, nil
	}

	if len(args) == 0 {
		return nil, nil // Empty resources will monitor all
	}

	resources := make([]Resource, 0, len(args))
	validResources := map[string]Resource{
		"cpu":  CPU,
		"disk": Disk,
		"mem":  Mem,
	}

	for _, arg := range args {
		r, ok := validResources[strings.ToLower(arg)]
		if !ok {
			return nil, fmt.Errorf("invalid resource: %s (valid: cpu, disk, mem)", arg)
		}
		resources = append(resources, r)
	}

	return resources, nil
}
