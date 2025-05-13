package cmd

import (
	"fmt"
	"os"

	"github.com/oliashish/vofo/logger"
	"github.com/spf13/cobra"
)

// Logger is initialized in init()
var log *logger.Logger

// State defines the agent execution modes.
type State string

const (
	All     State = "all"
	Monitor State = "monitor"
	Heal    State = "heal"
)

// Resource defines the monitorable resources.
type Resource string

const (
	CPU  Resource = "cpu"
	Disk Resource = "disk"
	Mem  Resource = "mem"
)

// RootCmd is the main CLI command.
var RootCmd = &cobra.Command{
	Use:   "vofo",
	Short: "A lightweight, cloud-agnostic monitoring and auto-healing tool",
	Long: `Vofo is a lightweight, cloud-agnostic tool for monitoring system resources
and performing auto-healing actions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// InitializeCmd handles project initialization.
var InitializeCmd = &cobra.Command{
	Use:     "init [config-path]",
	Short:   "Initialize project with a config file",
	Aliases: []string{"initialize"},
	Example: "vofo init config/config.json",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info(fmt.Sprintf("Running init command with config: %s", args[0]))
		return InitVofo(args[0])
	},
}

// StartCmd handles starting the agent.
var StartCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start the agent with specified modules",
	Aliases: []string{"s"},
	Example: `vofo start --all
vofo start --monitor
vofo start --heal`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var state State
		switch {
		case allFlag:
			state = All
		case monitorFlag:
			state = Monitor
		case healFlag:
			state = Heal
		default:
			state = Monitor
		}
		log.Info(fmt.Sprintf("Running start command with state: %s", state))
		return startAgent(state)
	},
}

// MonitorCmd handles monitoring specific resources.
var MonitorCmd = &cobra.Command{
	Use:   "monitor [resources...]",
	Short: "Monitor system resources (cpu, disk, mem) or all if none specified",
	Example: `vofo monitor cpu
vofo monitor cpu disk
vofo monitor --all`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		resources, err := parseResources(args, allResourcesFlag)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to parse resources: %s", err))
			return fmt.Errorf("invalid resources: %w", err)
		}
		log.Info(fmt.Sprintf("Running monitor command with resources: %v", resources))
		return startMonitor(resources)
	},
}

// Flags for start command.
var (
	allFlag     bool
	monitorFlag bool
	healFlag    bool
)

// Flags for monitor command.
var allResourcesFlag bool

// init sets up the CLI commands, flags, and logger.
func Init() {
	// Initialize logger
	var err error
	log, err = logger.NewLogger()
	if err != nil {
		_, _ = os.Stderr.WriteString("Failed to initialize logger: " + err.Error() + "\n")
		os.Exit(1)
	}

	// Add subcommands to root.
	RootCmd.AddCommand(InitializeCmd)
	RootCmd.AddCommand(StartCmd)
	RootCmd.AddCommand(MonitorCmd)

	// Define flags for start command.
	StartCmd.Flags().BoolVarP(&allFlag, "all", "A", false, "Start all modules")
	StartCmd.Flags().BoolVarP(&monitorFlag, "monitor", "M", false, "Start monitor module")
	StartCmd.Flags().BoolVarP(&healFlag, "heal", "H", false, "Start heal module")
	// Ensure flags are mutually exclusive.
	StartCmd.MarkFlagsMutuallyExclusive("all", "monitor", "heal")

	// Define flags for monitor command.
	MonitorCmd.Flags().BoolVarP(&allResourcesFlag, "all", "A", false, "Monitor all resources (cpu, disk, mem)")
}

// Execute runs the CLI.
func Execute() error {
	return RootCmd.Execute()
}
