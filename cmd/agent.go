package cmd

import (
	"fmt"
	"os"

	"github.com/oliashish/vofo/internals/monitor"
	"github.com/oliashish/vofo/logger"
	"github.com/spf13/cobra"
)

var log = logger.Logger()

type State string

const (
	All     State = "All"
	Monitor State = "Monitor"
	Heal    State = "Heal"
)

var rootCmd = &cobra.Command{
	Use:   "vofo",
	Short: "A lightweight, cloud agnostic monitoring and auto-healing tool",
	Long:  `A lightweight, cloud agnostic monitoring and auto-healing tool`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var initialize = &cobra.Command{
	Use:     "initialize",
	Short:   "Initlize project with a path to config file",
	Aliases: []string{"init"},
	Example: "vofo init path/to/config.json",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		InitializeVofo(args[0])
	},
}

var start = &cobra.Command{
	Use:     "start",
	Short:   "Start agent with configurable args",
	Aliases: []string{"s"},
	Example: `vofo start <defaults to monitor>
To use all modules please use <vofo start -A>`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case allFlag:
			startAgent(All)

		case monitFlag:
			startAgent(Monitor)

		case healFlag:
			startAgent(Heal)

		default:
			startAgent(Monitor)
		}
	},
}

var allFlag bool
var monitFlag bool
var healFlag bool

func Init() {
	rootCmd.AddCommand(initialize)

	rootCmd.AddCommand(start)

	start.Flags().BoolVarP(&allFlag, "all", "A", false, "Start all modules!!!")
	start.Flags().BoolVarP(&monitFlag, "monitor", "M", false, "Start monitor module!!!")
	start.Flags().BoolVarP(&healFlag, "heal", "H", false, "Start heal modules!!!")

}

func Execute() {
	logger := logger.Logger()
	if err := rootCmd.Execute(); err != nil {
		logger.Error(fmt.Sprintf("Error while executing Project: %s\n", err))
		os.Exit(1)
	}
}

func startAgent(s State) {
	logger := logger.Logger()
	switch s {
	case All:
		logger.Info("Start All Agent")
		monitor.CPU()
	case Monitor:
		logger.Info("Start Monitoring Agent")
		monitor.CPU()
	case Heal:
		logger.Info("Start Healing Agent")
	default:
		logger.Info("Default, running monitoring agent")
	}
}
