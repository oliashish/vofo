package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/oliashish/vofo/internals/config"
)

// InitializeVofo initializes the project with the given config path.
func InitVofo(configPath string) error {
	log.Info(fmt.Sprintf("Initializing project with config: %s", configPath))

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to read config file: %s", err))
		return fmt.Errorf("cannot read config: %w", err)
	}

	// Parse config
	var cfg config.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Error(fmt.Sprintf("Failed to parse cfg JSON: %s", err))
		return fmt.Errorf("invalid cfg format: %w", err)
	}

	// Validate cfg
	if cfg.CPUThreshold <= 0 || cfg.CPUThreshold > 100 {
		log.Error("Invalid CPU threshold: must be between 0 and 100")
		return fmt.Errorf("invalid CPU threshold: %.2f", cfg.CPUThreshold)
	}
	if cfg.RAMThreshold <= 0 || cfg.RAMThreshold > 100 {
		log.Error("Invalid RAM threshold: must be between 0 and 100")
		return fmt.Errorf("invalid RAM threshold: %.2f", cfg.RAMThreshold)
	}
	if cfg.DiskThreshold <= 0 || cfg.DiskThreshold > 100 {
		log.Error("Invalid Disk threshold: must be between 0 and 100")
		return fmt.Errorf("invalid Disk threshold: %.2f", cfg.DiskThreshold)
	}

	config.SetConfig(&cfg)

	log.Info("Configuration loaded successfully")
	return nil
}
