package worker

import (
	"fmt"
	"log"

	"github.com/fgrehm/sls-web/src/core/solver"
)

type ExperimentWorker struct {
	solver    solver.Service
	logPrefix string
}

func NewExperimentWorker(id string, solver solver.Service) *ExperimentWorker {
	return &ExperimentWorker{
		solver:    solver,
		logPrefix: fmt.Sprintf("[%s]", id),
	}
}

func (w *ExperimentWorker) Run(queue chan solver.ExperimentID, done chan struct{}) {
	for experimentID := range queue {
		w.Process(experimentID)
	}
	log.Printf("%s Shutting down...", w.logPrefix)
	done <- struct{}{}
}

func (w *ExperimentWorker) Process(experimentID solver.ExperimentID) {
	log.Printf("%s Processing experimentID %q", w.logPrefix, experimentID)

	err := w.solver.ProcessExperiment(experimentID)
	if err != nil {
		log.Printf("%s Error while processing %q: %s", w.logPrefix, experimentID, err)
	} else {
		log.Printf("%s Done with %q", w.logPrefix, experimentID)
	}
}
