package worker

import (
	"fmt"
	"log"

	"github.com/fgrehm/sls-web/src/core/solver"
)

type SolutionWorker struct {
	solver    solver.Service
	logPrefix string
}

func NewSolutionWorker(id string, solver solver.Service) *SolutionWorker {
	return &SolutionWorker{
		solver:    solver,
		logPrefix: fmt.Sprintf("[%s]", id),
	}
}

func (w *SolutionWorker) Run(queue chan solver.SolutionID, done chan struct{}) {
	for solutionID := range queue {
		w.Process(solutionID)
	}
	log.Printf("%s Shutting down...", w.logPrefix)
	done <- struct{}{}
}

func (w *SolutionWorker) Process(solutionID solver.SolutionID) {
	log.Printf("%s Processing solution %q", w.logPrefix, solutionID)

	err := w.solver.Solve(solutionID)
	if err != nil {
		log.Printf("%s Error while solving %q: %s", w.logPrefix, solutionID, err)
	} else {
		solution, err := w.solver.FindSolution(solutionID, nil)
		if err == nil {
			msg := "FOUND"
			if *solution.Errored {
				msg = "ERRORED"
			} else if !*solution.Found {
				msg = "NOT FOUND"
			}
			log.Printf("%s Done with %q: %s", w.logPrefix, solutionID, msg)
		}
	}
}
