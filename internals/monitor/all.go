package monitor

import (
	"fmt"
	"sync"

	"github.com/oliashish/vofo/internals/config"
	"github.com/oliashish/vofo/logger"
)

// All monitors all resources (CPU, RAM, Disk) concurrently.
func All() error {
	log, err := logger.NewLogger()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	cfg := config.GetConfig()
	if cfg == nil {
		log.Error("Configuration not initialized")
		return fmt.Errorf("configuration not initialized")
	}

	log.Info("Starting monitoring of all resources")

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	// Start CPU monitoring
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := CPU(); err != nil {
			errChan <- fmt.Errorf("cpu monitoring: %w", err)
		}
	}()

	// Start RAM monitoring
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := Mem(); err != nil {
			errChan <- fmt.Errorf("ram monitoring: %w", err)
		}
	}()

	// Start Disk monitoring
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := Disk(); err != nil {
			errChan <- fmt.Errorf("disk monitoring: %w", err)
		}
	}()

	// Wait for all goroutines to finish or an error to occur
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Return first error (if any)
	for err := range errChan {
		return err
	}

	return nil
}

// Heal performs healing actions (placeholder).
func Heal() error {
	log, err := logger.NewLogger()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	cfg := config.GetConfig()
	if cfg == nil {
		log.Error("Configuration not initialized")
		return fmt.Errorf("configuration not initialized")
	}

	if cfg.ServiceName == "" {
		log.Error("No service name configured for healing")
		return fmt.Errorf("no service name configured")
	}

	log.Info(fmt.Sprintf("Performing healing action: restarting service %s", cfg.ServiceName))
	// TODO: Healing
	// Placeholder: Implement healing (e.g., systemctl restart)
	return nil
}
