package main

import (
	"github.com/fgrehm/sls-web/src/worker"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(workCmd)
	workCmd.Flags().IntVarP(&SolverThreads, "threads-solver", "", SolverThreads, "Numbers of threads for solving models")
	workCmd.Flags().IntVarP(&ExperimenterThreads, "threads-experimenter", "", ExperimenterThreads, "Number of threads to start for scheduling experiments")
}

var workCmd = &cobra.Command{
	Use:   "work",
	Short: "Starts the worker process",

	Run: func(cmd *cobra.Command, args []string) {
		env := buildEnv()
		defer env.Shutdown()

		worker.Run(worker.RunOpts{
			Solver:           env.SolverSvc,
			AmqpChannel:      env.AmqpChannel,
			ExperimentsQueue: env.ExperimentsQueue,
			SolutionsQueue:   env.SolutionsQueue,
		})
	},
}
