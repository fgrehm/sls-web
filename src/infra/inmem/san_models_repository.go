package inmem

import (
	"sync"

	"github.com/fgrehm/sls-web/src/core/models"
)

type sanModelRepository struct {
	mtx        sync.RWMutex
	modelsById map[models.SanModelID]*models.SanModel
	models     []*models.SanModel
}

func (r *sanModelRepository) Store(m *models.SanModel) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.modelsById[m.ID]; !ok {
		r.models = append(r.models, m)
	}
	r.modelsById[m.ID] = m

	return nil
}

func (r *sanModelRepository) Find(id models.SanModelID) (*models.SanModel, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if val, ok := r.modelsById[id]; ok {
		return val, nil
	}
	return nil, models.ErrUnknown
}

func (r *sanModelRepository) All() ([]*models.SanModel, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	m := make([]*models.SanModel, 0, len(r.models))
	for _, val := range r.models {
		m = append(m, val)
	}
	return m, nil
}

// NewSanModelsRepository returns a new instance of a in-memory models repository.
func NewSanModelsRepository() models.Repository {
	return &sanModelRepository{
		modelsById: make(map[models.SanModelID]*models.SanModel),
		models:     []*models.SanModel{},
	}
}
