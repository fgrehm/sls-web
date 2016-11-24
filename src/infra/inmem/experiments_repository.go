package inmem

import (
	"sync"

	"github.com/fgrehm/sls-web/src/core/solver"
)

type experimentsRepository struct {
	mtx             sync.RWMutex
	experimentsById map[solver.ExperimentID]*solver.Experiment
	experiments     []*solver.Experiment
}

func (r *experimentsRepository) Store(experiment *solver.Experiment) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.experimentsById[experiment.ID]; !ok {
		r.experiments = append(r.experiments, experiment)
	}
	r.experimentsById[experiment.ID] = experiment

	return nil
}

func (r *experimentsRepository) Find(id solver.ExperimentID, externalKey *string) (*solver.Experiment, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if val, ok := r.experimentsById[id]; ok {
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

// NewExperimentsRepository returns a new instance of a in-memory experiment repository.
func NewExperimentsRepository() solver.ExperimentRepository {
	return &experimentsRepository{
		experimentsById: make(map[solver.ExperimentID]*solver.Experiment),
		experiments:     []*solver.Experiment{},
	}
}
