package main

import (
	"github.com/fgrehm/sls-web/src/core/compiler"
	"github.com/fgrehm/sls-web/src/core/models"
	"github.com/fgrehm/sls-web/src/core/renderer"
	"github.com/fgrehm/sls-web/src/core/solver"
	"github.com/fgrehm/sls-web/src/infra/inmem"
	"github.com/fgrehm/sls-web/src/infra/mongo"
	"github.com/fgrehm/sls-web/src/infra/rabbit"

	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2"
)

// REFACTOR: All of this "environment concept" needs to be revisited

type Env struct {
	CompilerSvc      compiler.Service
	RendererSvc      renderer.Service
	SolverSvc        solver.Service
	ModelsSvc        models.Service
	AmqpChannel      *amqp.Channel
	ExperimentsQueue amqp.Queue
	SolutionsQueue   amqp.Queue
	shutdown         func()
}

func (e *Env) Shutdown() {
	if e.shutdown != nil {
		e.shutdown()
	}
}

func buildEnv() *Env {
	if DemoMode {
		return buildDemoModeEnv()
	} else if Standalone {
		return buildStandaloneEnv()
	}

	return buildDefaultEnv()
}

func buildDemoModeEnv() *Env {
	scheduler := inmem.NewScheduler()

	solutionsRepo := inmem.NewSolutionsRepository()
	experimentsRepo := inmem.NewExperimentsRepository()
	modelsRepo := inmem.NewSanModelsRepository()

	compilerSvc := compiler.NewService()
	rendererSvc := renderer.NewService(compilerSvc)
	solverSvc := solver.NewService(solutionsRepo, experimentsRepo, scheduler, compilerSvc)
	modelsSvc := models.NewService(modelsRepo, compilerSvc, rendererSvc, solverSvc)

	scheduler.StartSolutionsWorker(solverSvc, SolverThreads)
	scheduler.StartExperimentsWorker(solverSvc, ExperimenterThreads)

	return &Env{
		CompilerSvc: compilerSvc,
		RendererSvc: rendererSvc,
		SolverSvc:   solverSvc,
		ModelsSvc:   modelsSvc,
	}
}

func buildStandaloneEnv() *Env {
	scheduler := inmem.NewScheduler()

	mongoSession, closeMongo := prepareMongoSession()
	modelsRepo := mongo.NewSanModelsRepository(mongoSession, MongoDatabase)
	solutionsRepo := mongo.NewSolutionsRepository(mongoSession, MongoDatabase)
	experimentsRepo := mongo.NewExperimentsRepository(mongoSession, MongoDatabase, solutionsRepo)

	compilerSvc := compiler.NewService()
	rendererSvc := renderer.NewService(compilerSvc)
	solverSvc := solver.NewService(solutionsRepo, experimentsRepo, scheduler, compilerSvc)
	modelsSvc := models.NewService(modelsRepo, compilerSvc, rendererSvc, solverSvc)

	scheduler.StartSolutionsWorker(solverSvc, SolverThreads)
	scheduler.StartExperimentsWorker(solverSvc, ExperimenterThreads)

	return &Env{
		CompilerSvc: compilerSvc,
		RendererSvc: rendererSvc,
		SolverSvc:   solverSvc,
		ModelsSvc:   modelsSvc,
		shutdown: func() {
			closeMongo()
			scheduler.Shutdown()
		},
	}
}

func buildDefaultEnv() *Env {
	amqpConn, err := amqp.Dial(RabbitURL)
	failOnError(err, "Failed to connect to RabbitMQ")

	amqpChannel, err := amqpConn.Channel()
	failOnError(err, "Failed to open a channel")

	solutionsQueue, err := amqpChannel.QueueDeclare(
		"solutions", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	experimentsQueue, err := amqpChannel.QueueDeclare(
		"experiments", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	scheduler := rabbit.NewScheduler(amqpChannel, experimentsQueue, solutionsQueue)

	mongoSession, closeMongo := prepareMongoSession()
	modelsRepo := mongo.NewSanModelsRepository(mongoSession, MongoDatabase)
	solutionsRepo := mongo.NewSolutionsRepository(mongoSession, MongoDatabase)
	experimentsRepo := mongo.NewExperimentsRepository(mongoSession, MongoDatabase, solutionsRepo)

	compilerSvc := compiler.NewService()
	rendererSvc := renderer.NewService(compilerSvc)
	solverSvc := solver.NewService(solutionsRepo, experimentsRepo, scheduler, compilerSvc)
	modelsSvc := models.NewService(modelsRepo, compilerSvc, rendererSvc, solverSvc)

	return &Env{
		CompilerSvc:      compilerSvc,
		RendererSvc:      rendererSvc,
		SolverSvc:        solverSvc,
		ModelsSvc:        modelsSvc,
		AmqpChannel:      amqpChannel,
		ExperimentsQueue: experimentsQueue,
		SolutionsQueue:   solutionsQueue,

		shutdown: func() {
			closeMongo()
			amqpChannel.Close()
			amqpConn.Close()
		},
	}
}

func prepareMongoSession() (*mgo.Session, func()) {
	mongoSession, err := mgo.Dial(MongoURL)
	failOnError(err, "Failed to connect to mongo")
	// Optional. Switch the mongoSession to a monotonic behavior.
	mongoSession.SetMode(mgo.Monotonic, true)

	return mongoSession, func() {
		mongoSession.Close()
	}
}
