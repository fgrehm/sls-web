package main

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().IntVarP(&SolverThreads, "threads-solver", "", SolverThreads, "Numbers of threads for solving models")
	startCmd.Flags().IntVarP(&ExperimenterThreads, "threads-experimenter", "", ExperimenterThreads, "Number of threads to start for scheduling experiments")
	startCmd.Flags().BoolVarP(&DemoMode, "demo", "", DemoMode, "Starts SLS Web in demo mode where all data is kept in memory and discarded on shutdown or an error")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the app in the standalone mode",

	Run: func(cmd *cobra.Command, args []string) {
		Standalone = true
		serveCmd.Run(cmd, args)
	},
}
