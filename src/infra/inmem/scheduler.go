package inmem

import (
	"fmt"
	"log"

	"github.com/fgrehm/sls-web/src/core/solver"
	"github.com/fgrehm/sls-web/src/worker"
)

type InMemScheduler struct {
	solutionQueue     chan solver.SolutionID
	solutionWorkers   []*worker.SolutionWorker
	experimentQueue   chan solver.ExperimentID
	experimentWorkers []*worker.ExperimentWorker
	done              chan struct{}
}

func (s *InMemScheduler) ScheduleSolution(solution *solver.Solution) error {
	s.solutionQueue <- solution.ID
	return nil
}

func (s *InMemScheduler) ScheduleExperiment(experiment *solver.Experiment) error {
	s.experimentQueue <- experiment.ID
	return nil
}

func (s *InMemScheduler) StartSolutionsWorker(svc solver.Service, n int) {
	log.Printf("Starting %d solution workers", n)
	for i := 0; i < n; i++ {
		worker := worker.NewSolutionWorker(fmt.Sprintf("solutions-%d", i+1), svc)
		s.solutionWorkers = append(s.solutionWorkers, worker)
		go worker.Run(s.solutionQueue, s.done)
	}
}

func (s *InMemScheduler) StartExperimentsWorker(svc solver.Service, n int) {
	log.Printf("Starting %d experiment workers", n)
	for i := 0; i < n; i++ {
		worker := worker.NewExperimentWorker(fmt.Sprintf("experiments-%d", i+1), svc)
		s.experimentWorkers = append(s.experimentWorkers, worker)
		go worker.Run(s.experimentQueue, s.done)
	}
}

func (s *InMemScheduler) Shutdown() {
	log.Println("Shutting down background workers...")
	close(s.experimentQueue)
	close(s.solutionQueue)
	totalWorkers := len(s.solutionWorkers) + len(s.experimentWorkers)
	for i := 0; i < totalWorkers; i++ {
		<-s.done
	}
}

func NewScheduler() *InMemScheduler {
	return &InMemScheduler{
		done:              make(chan struct{}),
		solutionQueue:     make(chan solver.SolutionID, 20),
		solutionWorkers:   []*worker.SolutionWorker{},
		experimentQueue:   make(chan solver.ExperimentID, 10),
		experimentWorkers: []*worker.ExperimentWorker{},
	}
}
