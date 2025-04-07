package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/oliashish/vofo/config"
)

func InitializeVofo(path string) {
	configFile, err := os.ReadFile(path)

	if err != nil {
		log.Error(fmt.Sprintf("Error while executing Project: %s\n", err))
		os.Exit(1)
	}

	var userConfig config.Config
	parseErr := json.Unmarshal(configFile, &userConfig)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to parse your configs from JSON file at : %s ERROR:  %s\n", configFile, parseErr))
		os.Exit(1)
	}

	log.Info("Please check your config settings for the machine and confirm")
	log.Info("Only [YES] is allowed: ")

	var userIn string
	fmt.Scan(&userIn)

	if err != nil {
		log.Error(fmt.Sprintf("Only [YES] is allowed: %s\n", err))
		os.Exit(1)
	}

	if strings.ToLower(userIn) != "yes" {
		log.Error("Only [YES] is allowed. Exiting...")
		os.Exit(1)
	}

	log.Info("initializing your machine....")

}
