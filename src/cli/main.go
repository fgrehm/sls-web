package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "sls-web",
	Short: "SAN Lite-Solver on the Web",
}

var (
	MongoURL            = "localhost:27017"
	MongoDatabase       = "slsweb"
	RabbitURL           = "amqp://guest:guest@localhost:5672/"
	SolverThreads       = 1
	ExperimenterThreads = 1
	DemoMode            = false
	Standalone          = false
)

func main() {
	if os.Getenv("MONGO_URL") != "" {
		MongoURL = os.Getenv("MONGO_URL")
	} else if os.Getenv("MONGODB_URI") != "" {
		MongoURL = os.Getenv("MONGODB_URI")
		MongoDatabase = ""
	}

	if os.Getenv("RABBIT_URL") != "" {
		RabbitURL = os.Getenv("RABBIT_URL")
	} else if os.Getenv("CLOUDAMQP_URL") != "" {
		RabbitURL = os.Getenv("CLOUDAMQP_URL")
	}

	RootCmd.PersistentFlags().StringVarP(&MongoURL, "mongo-url", "m", MongoURL, "URL for the MongoDB instance to conect to")
	RootCmd.PersistentFlags().StringVarP(&MongoDatabase, "mongo-db", "d", MongoDatabase, "MongoDB database to use")
	RootCmd.PersistentFlags().StringVarP(&RabbitURL, "rabbit-url", "r", RabbitURL, "URL for the RabbitMQ instance to connect to")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
