package worker

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/fgrehm/sls-web/src/core/solver"

	"github.com/streadway/amqp"
)

type RunOpts struct {
	Solver           solver.Service
	AmqpChannel      *amqp.Channel
	ExperimentsQueue amqp.Queue
	SolutionsQueue   amqp.Queue
}

func Run(runOpts RunOpts) {
	svc := runOpts.Solver
	ch := runOpts.AmqpChannel
	experimentsQueue := runOpts.ExperimentsQueue
	solutionsQueue := runOpts.SolutionsQueue

	experimentsMsgs, err := ch.Consume(
		experimentsQueue.Name, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	failOnError(err, "Failed to register an experiment consumer")

	solutionsMsgs, err := ch.Consume(
		solutionsQueue.Name, // queue
		"",                  // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	failOnError(err, "Failed to register a solutions consumer")

	done := make(chan struct{})

	solutionWorker := NewSolutionWorker("solutions", svc)
	go func() {
		for msg := range solutionsMsgs {
			log.Printf("Received a message on SOLUTIONS queue: %s", msg.Body)
			solutionWorker.Process(solver.SolutionID(msg.Body))
		}
		log.Printf("Shutting down solution worker thread")
		done <- struct{}{}
	}()

	experimentWorker := NewExperimentWorker("experiments", svc)
	go func() {
		for msg := range experimentsMsgs {
			log.Printf("Received a message on EXPERIMENTS queue: %s", msg.Body)
			experimentWorker.Process(solver.ExperimentID(msg.Body))
		}
		log.Printf("Shutting down experiment worker thread")
		done <- struct{}{}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	<-c

	log.Printf("Shutting down workers...")

	ch.Close()
	<-done
	<-done

	log.Printf("Done")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
