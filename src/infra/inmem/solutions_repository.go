package inmem

import (
	"sync"

	"github.com/fgrehm/sls-web/src/core/solver"
)

type solutionsRepository struct {
	mtx           sync.RWMutex
	solutionsById map[solver.SolutionID]*solver.Solution
	solutions     []*solver.Solution
}

func (r *solutionsRepository) Store(solution *solver.Solution) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.solutionsById[solution.ID]; !ok {
		r.solutions = append(r.solutions, solution)
	}
	r.solutionsById[solution.ID] = solution

	return nil
}

func (r *solutionsRepository) Find(id solver.SolutionID, externalKey *string) (*solver.Solution, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if val, ok := r.solutionsById[id]; ok {
		if externalKey == nil {
			return val, nil
		}
		if val.ExternalKey == nil {
			return nil, solver.ErrUnknown
		}
		if *val.ExternalKey == *externalKey {
			return val, nil
		}
	}
	return nil, solver.ErrUnknown
}

func (r *solutionsRepository) FindForReuse(externalKey, modelHash string, maxIterations uint, tolerance float64) (*solver.Solution, error) {
	return nil, solver.ErrUnknown
}

// NewSolutionsRepository returns a new instance of a in-memory solution repository.
func NewSolutionsRepository() solver.SolutionRepository {
	return &solutionsRepository{
		solutionsById: make(map[solver.SolutionID]*solver.Solution),
		solutions:     []*solver.Solution{},
	}
}
