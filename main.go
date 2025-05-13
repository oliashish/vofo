package main

import (
	"os"

	"github.com/oliashish/vofo/cmd"
	"github.com/oliashish/vofo/logger"
)

// main initializes the application and handles CLI execution.
func main() {
	// Initialize logger
	log, err := logger.NewLogger()
	if err != nil {
		// Print to stderr since logger failed
		_, _ = os.Stderr.WriteString("Failed to initialize logger: " + err.Error() + "\n")
		os.Exit(1)
	}

	// Set up logger for debug and info
	log.Debug("Setting up Vofo agent...")
	log.Info("Vofo agent starting...")

	// Initialize and execute CLI commands
	cmd.Init()
	if err := cmd.Execute(); err != nil {
		log.Error("Failed to execute CLI: " + err.Error())
		os.Exit(1)
	}
}
