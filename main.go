package main

import (
	"github.com/oliashish/vofo/cmd"
	"github.com/oliashish/vofo/logger"
)

// TODO: handle any error while creating zap logger

func main() {
	logger := logger.Logger()

	logger.Debug("Setting up Agent...")
	logger.Info("Setting up Agent...")
	cmd.Init()
	cmd.Execute()

}
