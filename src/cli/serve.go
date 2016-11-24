package main

import (
	"os"

	"github.com/fgrehm/sls-web/src/api"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the Web server",

	Run: func(cmd *cobra.Command, args []string) {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		env := buildEnv()
		defer env.Shutdown()

		api.Run(api.RunOpts{
			Port:         port,
			CompilerSvc:  env.CompilerSvc,
			RendererSvc:  env.RendererSvc,
			SanModelsSvc: env.ModelsSvc,
			SolverSvc:    env.SolverSvc,
		})
	},
}
